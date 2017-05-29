package main
import (
	"fmt"
    "log"
    "os"
    "net/http"
    "encoding/base64"
    "crypto/rand"
    qrcode "github.com/skip2/go-qrcode"
    btckey "github.com/vsergeev/btckeygenie/btckey"
)

const br = "<br>"

func check(err error) {
    if err != nil {
        panic(err)
    }
}

func byteString(b []byte) (s string) {
	s = ""
	for i := 0; i < len(b); i++ {
		s += fmt.Sprintf("%02X", b[i])
	}
	return s
}

func generateQRCode(input string) (string) {
    var png []byte
    png, err := qrcode.Encode(input, qrcode.Medium, 128)
    check(err)
    imgBase64Str := base64.StdEncoding.EncodeToString(png)
    return "<img src=\"data:image/png;base64," + imgBase64Str + "\" />"
}

func main() {
	fmt.Println("Running local server @ http://localhost:" + os.Getenv("PORT"))
    
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        url := r.FormValue("url")
        if url == "" {
            var priv btckey.PrivateKey
            var err error
            priv, err = btckey.GenerateKey(rand.Reader)
            check(err)
        
            address_compressed := priv.ToAddress()
            pub_bytes_compressed := priv.PublicKey.ToBytes()
            pub_bytes_compressed_str := byteString(pub_bytes_compressed)
            address_uncompressed := priv.ToAddressUncompressed()
            pub_bytes_uncompressed := priv.PublicKey.ToBytesUncompressed()
            pub_bytes_uncompressed_str := byteString(pub_bytes_uncompressed)
            wif := priv.ToWIF()
            pri_bytes := priv.ToBytes()
            pri_bytes_str := byteString(pri_bytes)
            
            w.Header().Set("Content-Type", "text/html; charset=utf-8")
            fmt.Fprintf(w, "<html><head></head><body><div style='width: 750px; word-wrap: break-word;'>")
            
            fmt.Fprintf(w, "<strong>Bitcoin Address:</strong>" + br + address_uncompressed + br)
            fmt.Fprintf(w, generateQRCode(address_uncompressed) + br)
            fmt.Fprintf(w, "<strong>Bitcoin Address (Compressed):</strong>" + br + address_compressed + br)
            fmt.Fprintf(w, generateQRCode(address_compressed) + br + br)
            
            fmt.Fprintf(w, "<strong>Public Key:</strong>" + br + pub_bytes_uncompressed_str + br)
            fmt.Fprintf(w, "<strong>Public Key (Compressed):</strong>" + br + pub_bytes_compressed_str + br + br + br + br)
            
            fmt.Fprintf(w, "<strong>Private Key:</strong>" + br + pri_bytes_str + br)
            fmt.Fprintf(w, generateQRCode(pri_bytes_str) + br)
            fmt.Fprintf(w, "<strong>Private Key (WIF):</strong>" + br + wif + br)
            fmt.Fprintf(w, generateQRCode(wif) + br)
            fmt.Fprintf(w, "</div></body></html>")
            return
        }
    })

    log.Fatal(http.ListenAndServe(":" + os.Getenv("PORT"), nil))
}
