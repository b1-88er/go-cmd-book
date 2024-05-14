#  Viper

```bash
PSCAN_HOSTS_FILE=newFile.hosts ./pScan hosts l
PSCAN_HOSTS_FILE=newFile.hosts ./pScan hosts add host01 host02
PSCAN_HOSTS_FILE=newFile.hosts ./pScan hosts list
PSCAN_HOSTS_FILE=newFile.hosts ./pScan hosts l
echo "hosts-file: newFile.hosts" > config.yaml
./pScan hosts l --config config.yaml
```