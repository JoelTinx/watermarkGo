# WatermarkGO: #

----------

API Go para colocar en una imagen (JPG) un sello de aguas (otra imagen en formato PNG).

# Instalación

go get github.com/JoelTinx/watermarkGo

```go

import w "github.com/JoelTinx/watermarkGo"

func main() {
	// 1: Para poner sello de agua a una sola imagen
	w.SetImageWaterMark("imagen.jpg", "watermark.png")

	// 2: Para poner sello de agua a todo el directorio "images"
	w.SetDirWaterMark("images", "watermark.png")

    //Nota: Se debe configurar el directorio de salida, por defecto se cuardará en "output".
}
```

# Autor
Sígueme en twitter: [Joel Tinx](https://twitter.com/joeltinx "https://twitter.com/joeltinx")
