# rucpy-downloader

Downloader written in golang to download the public data files(RUC Paraguay) from [set.gov.py](https://www.set.gov.py/portal/PARAGUAY-SET/InformesPeriodicos?folder-id=repository:collaboration:/sites/PARAGUAY-SET/categories/SET/Informes%20Periodicos/listado-de-ruc-con-sus-equivalencias).

The downloader will download the public data files (RUC) and extract them into the specified destination directory.

## Usage

<br/>

### From [source](https://github.com/bitebait/rucpy-downloader.git)

```bash
git clone https://github.com/bitebait/rucpy-downloader.git
cd rucpy-downloader/
go run main.go -d "destination directory"
```

### From [release](https://github.com/bitebait/rucpy-downloader/releases)

<br/>

#### Windows (CMD)

```bash
rucdownloader.exe -d "destination directory"
```

#### Linux (terminal)

```bash
./rucdownloader -d "destination directory"
```

## Contributing

1. Fork it!
2. Create your feature branch: `git checkout -b my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin my-new-feature`
5. Submit a pull request :D

## License

This project is licensed under the terms of the MIT License - see the [LICENSE](LICENSE) file for details
