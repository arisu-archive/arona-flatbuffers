import importlib
from importlib import metadata
import json
from pathlib import Path
import unittest


ROOT = Path(__file__).resolve().parents[2]
MANIFEST_PATH = ROOT / ".release-please-manifest.json"
RELEASE_CONFIG_PATH = ROOT / "release-please-config.json"
LOCK_PATH = ROOT / "uv.lock"
DISTRIBUTION_NAME = "arona-flatbuffers"


class DistributionContractTests(unittest.TestCase):
    def import_or_fail(self, module_name: str):
        try:
            return importlib.import_module(module_name)
        except ModuleNotFoundError as error:
            self.fail(f"installed distribution is missing {module_name}: {error}")

    def test_installed_version_matches_release_manifest(self):
        expected = json.loads(MANIFEST_PATH.read_text(encoding="utf-8"))["."]
        try:
            installed = metadata.version(DISTRIBUTION_NAME)
        except metadata.PackageNotFoundError:
            self.fail(f"distribution is not installed: {DISTRIBUTION_NAME}")

        self.assertEqual(installed, expected)

    def test_release_please_updates_lockfile_version(self):
        config = json.loads(RELEASE_CONFIG_PATH.read_text(encoding="utf-8"))
        extra_files = config["packages"]["."]["extra-files"]
        expected_updater = {
            "type": "toml",
            "path": "uv.lock",
            "jsonpath": "$.package[0].version",
        }
        lock_text = LOCK_PATH.read_text(encoding="utf-8")
        first_package = lock_text.split("[[package]]", maxsplit=1)[1]
        first_package_record = first_package.split("[[package]]", maxsplit=1)[0]

        self.assertIn(expected_updater, extra_files)
        self.assertIn(f'name = "{DISTRIBUTION_NAME}"', first_package_record)

    def test_flatdata_object_api_round_trips(self):
        module = self.import_or_fail("FlatData.BlendInfo")
        original = module.BlendInfoT(from_=11, to=29, blend=0.5)

        restored = module.BlendInfoT.from_bytes(original.to_bytes())

        self.assertEqual(restored.from_, 11)
        self.assertEqual(restored.to, 29)
        self.assertAlmostEqual(restored.blend, 0.5)

    def test_mx_namespace_is_installed(self):
        module = self.import_or_fail("MX.Data.Excel.WorldRaidStageRewardExcel")

        self.assertTrue(hasattr(module, "WorldRaidStageRewardExcelT"))


if __name__ == "__main__":
    unittest.main()
