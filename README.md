# OTH

##### Odoo To Hexya
translate an Odoo module into a Hexya module at best as it can

### Usage
go to the 'Launcher' directory.

Compile the launcher using `go build`
and start the server with the command `./Launcher <odoo_dir> <output_dir>`

e.g.

```sh
cd Launcher
./Launcher ../input/mail ../output/mail
```

multiple prompts should show up on the console, asking you the path of the go modules corresponding to the dependencies.
enter "skip" if you wish to ignore this dependency. leave empty if you want to search in the default hexya-addons repositories

```
HexTranslate - 2019-06-25 10:54:53 - [Warn]:    Some dependencies found are not recognized.
                please specify for each one its hexya module path
                Leave empty for default hexya-addons folder
                Exemple: web -> github.com/hexya-addons/web


HexTranslate - 2019-06-25 10:54:53 - [Warn]:    please specify hexya module path for:   baseSetup
skip
HexTranslate - 2019-06-25 10:54:56 - [Warn]:    please specify hexya module path for:   bus

HexTranslate - 2019-06-25 10:54:57 - [Warn]:    please specify hexya module path for:   webTour
github.com/another/repository
HexTranslate - 2019-06-25 10:55:18 - [Info]:    These module paths will attempt to be loaded:
HexTranslate - 2019-06-25 10:55:18 - [Info]:            github.com/hexya-addons/base
HexTranslate - 2019-06-25 10:55:18 - [Info]:            github.com/hexya-addons/bus
HexTranslate - 2019-06-25 10:55:18 - [Info]:            github.com/another/repository
HexTranslate - 2019-06-25 10:55:18 - [Info]:    Confirm?   [Y/N]
y
```

After a while, the server should be started, and the web interface accessible at
`http://localhost:8080`

Once on the web interface, click "create"
then, in the new form that appears, change the fields according to what you want.
then, click on generate.

The new hexya module will be written where you specified (output path)

