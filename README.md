## CD Template
if you use github actions, use example:
```yaml
- uses: actions/setup-node@v2
- name: Get CD Tools
  run: |
    curl -O https://mss-boot-io.github.io/cd-template/latest/linux_amd64.tar.gz
    tar -zxvf linux_amd64.tar.gz 
    ./cd-template --config=[you config path]
```
OR
```yaml
- uses: actions/setup-node@v2
- name: Get CD Tools
  run: |
    curl -O https://mss-boot-io.github.io/cd-template/latest/linux_amd64.tar.gz
    tar -zxvf linux_amd64.tar.gz 
    ./cd-template --namespace=prod --app=cd --service=temp --image=werrerw:v1 --configDataFiles=./cfg.local.yaml --importEnvNames=STAGE
```
