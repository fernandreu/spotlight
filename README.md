# Windows Spotlight Fetcher

This simple Go program allows you to obtain all interesting pictures from Windows Spotlight (i.e. the lock screen), to
be used for any other purpose (e.g. to set them up as wallpapers).

#### Usage from command line

```
go run .\cmd\main.go [DestinationFolder]
```

If the destination folder is not provided, the files will be copied to current one.

#### Main usage from Go

```
import (
	"github.com/fernandreu/spotlight/pkg"
)

...

app.ProcessFiles(origin, destination)
```

Where `origin` will typically be `app.GetDefaultSpotlightFolder()`.

#### Example output

```
Copying all available Windows Spotlight pictures

Origin folder:
C:\Users\UserName\AppData\Local\Packages\Microsoft.Windows.ContentDeliveryManager_cw5n1h2txyewy\LocalState\Assets

Destination folder:
.\

File copied: 283f752828f08f0458835888361a313aae2b5a293b3f52dedf82032c10c2938a.jpeg (1920x1080, 853 kb)
Finished copying; 1 new pictures found in total
```

## How it works

Windows periodically generates new pictures under the following folder:

```
%USERPROFILE%\AppData\Local\Packages\Microsoft.Windows.ContentDeliveryManager_cw5n1h2txyewy\LocalState\Assets
```

The pictures used by Windows Spotlight will also be contained in there. However, there are a few caveats:

- The pictures do not contain any extension, and can be either jpg or png
- There may be smaller versions of the same pictures, such as thumbnails / banners
- A picture previously shown by Windows Spotlight might appear there again but with a different name
- Pictures periodically disappear from there

This program overcomes this limitations by:

1. Scanning all files in that folder
2. Comparing them with the ones in the destination folder (via MD5 checksums)
3. If not in the destination folder, checking if the picture is in landscape orientation and has a minimum resolution of
1024x768
4. If so, copying the file to the destination folder with the appropriate extension
5. Finally, writing all MD5 checksums into a `CheckSums.txt` file in the destination folder which makes the process
faster next time it is launched

A single run of the program per day will probably give you all the pictures you can possibly get using this method. 
Seeing a change in your screen lock picture can be a good reminder to run it again.
