curl -s https://raw.githubusercontent.com/asfu222/BACNLocalizationResources/refs/heads/main/ba.env -o ba.env
export $(grep -v '^#' ba.env | xargs)
echo Using catalog url: $ADDRESSABLE_CATALOG_URL
mkdir -p resources/TableBundles
curl "${ADDRESSABLE_CATALOG_URL}/TableBundles/Excel.zip" -o resources/TableBundles/Excel.zip
curl "${ADDRESSABLE_CATALOG_URL}/TableBundles/ExcelDB.db" -o resources/TableBundles/ExcelDB.db

git clone --depth=1 https://github.com/asfu222/BlueArchiveLocalizationTools
pip3 install -r BlueArchiveLocalizationTools/requirements.txt
python3 fetch_resources.py