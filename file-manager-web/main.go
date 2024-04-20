package main

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func formatSize(size int64) string {
	const gb = 1024 * 1024 * 1024
	if size >= gb {
		return fmt.Sprintf("%.2f GB", float64(size)/gb)
	}
	const mb = 1024 * 1024
	if size >= mb {
		return fmt.Sprintf("%.2f MB", float64(size)/mb)
	}
	return fmt.Sprintf("%.2f KB", float64(size)/1024)
}

func isMedia(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".webp", ".jpeg", ".wav", ".png", ".gif", ".mp4", ".mkv", ".ogg", ".flac", ".mp3", ".m4a", ".webm", ".opus":
		return true
	default:
		return false
	}
}

func serveFile(w http.ResponseWriter, r *http.Request, filePath string) {
	// Decode the file path sent by the client
	decodedFilePath, err := url.PathUnescape(filePath)
	if err != nil {
		log.Println("Error decoding file path:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	file, err := os.Open(decodedFilePath)
	if err != nil {
		log.Println("Error opening file:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		log.Println("Error getting file info:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileInfo.Name()))

	contentType := mime.TypeByExtension(filepath.Ext(fileInfo.Name()))
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	w.Header().Set("Content-Type", contentType)

	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	_, err = io.Copy(w, file)
	if err != nil {
		log.Println("Error streaming file content:", err)
		// Note: Do not call http.Error here, as headers have already been written
		return
	}
}

func readFileContent(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func ContainsString(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

// this for prevent people to download system files and access private folders
var skipName = []string{"desktop.ini", "System32", "SysWOW64", "Recovery", "thumbs.db", "$RECYCLE.BIN", "System Volume Information", "$WinREAgent", "hiberfil.sys", "pagefile.sys", "swapfile.sys", "Documents and Settings", "DumpStack.log.tmp", "$Recycle.Bin"}

func handler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("file")
	if filePath != "" {
		serveFile(w, r, filePath)
		return
	}

	dirPath := filepath.Clean((r.URL.Query().Get("dir")))
	if dirPath == "" {
		http.Redirect(w, r, "/?dir=/", http.StatusFound)
		return
	}

	var files []fs.DirEntry
	var err error

	// Decode only the "+", "-", "%2B", and "%2D" characters
	dirPath = strings.ReplaceAll(dirPath, "%2B", "+")
	dirPath = strings.ReplaceAll(dirPath, "%2D", "-")
	dirPath = strings.ReplaceAll(dirPath, "\\+", "+")
	dirPath = strings.ReplaceAll(dirPath, "\\-", "-")
	// // Check if the directory path contains spaces or special characters
	// if strings.ContainsAny(dirPath, " \"") {
	// 	// Directory path contains spaces or special characters, wrap it in double quotes
	// 	dirPath = `"` + dirPath + `"`
	// }

	files, err = os.ReadDir(dirPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type File struct {
		Name  string
		Size  string
		Path  string
		IsDir bool
	}

	var filesList []File

	if dirPath != "." {
		parentDirPath := filepath.Dir(dirPath)
		filesList = append(filesList, File{
			Name:  "Go Up",
			Size:  "-",
			Path:  fmt.Sprintf("/?dir=%s", filepath.ToSlash(parentDirPath)),
			IsDir: true,
		})
	}

	for _, file := range files {
		// for skipping file or directory that we don't want to show
		if (file.IsDir() || !file.IsDir()) && ContainsString(skipName, file.Name()) {
			continue
		}

		// for detected that it's a directory yeah kinda clever stuff
		var isDir bool
		if file.IsDir() {
			isDir = true
		}

		fileInfo, err := file.Info()
		if err != nil {
			fmt.Printf("Error getting file info: %v\n", err)
		}
		filePath := filepath.Join(dirPath, file.Name())
		filesList = append(filesList, File{
			Name: file.Name(),
			Size: func() string {
				if isDir {
					return ""
				}
				return formatSize(fileInfo.Size())
			}(),
			Path: fmt.Sprintf("/?%s=%s", (func() string {
				if file.IsDir() {
					return "dir"
				}
				return "file"
			})(), filepath.ToSlash(filePath)),
			IsDir: isDir,
		})
	}

	funcMap := template.FuncMap{
		"isMedia":         isMedia,
		"readFileContent": readFileContent,
		"hasSuffix": func(s, suffix string) bool {
			return strings.HasSuffix(s, suffix)
		},
		// blame golang template for this because if not is there
		// the linter will complains about why this function not exists
		// even though i make function on javascript
		"replacePlus": func(s string) template.URL {
			s = strings.ReplaceAll(s, "+", "%2B")
			s = strings.ReplaceAll(s, "#", "%23")
			s = strings.ReplaceAll(s, "[", "%5B")
			s = strings.ReplaceAll(s, "]", "%5D")
			return template.URL(s)
		},
	}
	tmpl := template.Must(template.New("files").Funcs(funcMap).Parse(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>Filemanager Web</title>
			<style>
				table {
					border-collapse: collapse;
					width: 100%;
				}
				th, td {
					border: 1px solid black;
					padding: 6px;
					text-align: left;
				}
				/* Add gap between Name and Size columns */
				td:nth-child(2) {
					padding-left: 3px;
				}
			</style>
				<script>
				function replaceSpecialCharacters(str) {
					return str.replace(/[#\[\]]/g, function(match) {
						switch (match) {
							case '#':
								return '%23';
							case '[':
								return '%5B';
							case ']':
								return '%5D';
							case '+':
								return '%2B';
							default:
								return match;
						}
					});
				}

			    // this for for handling where cause by browser who dosen't encode properly for something like "+"
				// so to reduce i'm being insane i wrote this 
				function replacePlusWithEncoded() {
					const links = document.querySelectorAll('a');
					links.forEach(link => {
						link.href = replaceSpecialCharacters(link.href);
					});
			
					const mediaElements = document.querySelectorAll('img, video, source, audio');
					mediaElements.forEach(element => {
						if (element.hasAttribute('src')) {
							element.src = replaceSpecialCharacters(element.src);
						}
					});
				}
			
				window.onload = function() {
					// to prevent where localhost:8080 will redirect you to path where this program is located 
					// so to prevent that i wrote this to redirect user to /?dir=/ when they visit localhost:8080
					if (window.location.pathname === "/" && window.location.search === "") {
						window.location.href = "/?dir=/";
					}
			
					replacePlusWithEncoded();
				};
			</script>
		</head>
		<body>
			<noscript>
        		<p>Please enable JavaScript in your browser to view this page correctly.</p>
    		</noscript>
			<h1>Files and Folders</h1>
			<table>
				<tr>
					<th>Name</th>
					<th>Size</th>
				</tr>
				{{range .}}
				<tr>
					<td>
						{{if .IsDir}}
							<a href="{{.Path}}">{{.Name}}</a>
						{{else}}
							{{if isMedia .Name}}
								{{if or (and (hasSuffix .Name ".jpg") (hasSuffix .Name ".jpeg")) (hasSuffix .Name ".png")}}
									<img src="{{.Path}}" data-src="{{.Path}}" alt="{{.Name}}" width="250"><br>
									<a href="{{.Path}}">{{.Name}}</a>
								{{else if or (hasSuffix .Name ".webp") (hasSuffix .Name ".gif")}}
									<img src="{{.Path}}" data-src="{{.Path}}" alt="{{.Name}}" width="250"><br>
									<a href="{{.Path}}">{{.Name}}</a>
								{{else if or (hasSuffix .Name ".mp4") (hasSuffix .Name ".webm")}}
									<video controls width="320" height="240" preload="none">
										<source src="{{.Path}}" type="video/mp4">
										<source src="{{.Path}}" type="video/webm">
										Your browser does not support the video tag.
									</video><br>
									<a href="{{.Path}}">{{.Name}}</a>
								{{else if hasSuffix .Name ".wav"}}
									<audio controls preload="none">
										<source src="{{.Path}}" type="audio/wav">
										Your browser does not support the audio tag.
									</audio><br>
									<a href="{{.Path}}">{{.Name}}</a>
								{{else if hasSuffix .Name ".m4a"}}
									<audio controls preload="none">
										<source src="{{.Path}}" type="audio/mpeg">
										Your browser does not support the audio tag.
									</audio><br>
									<a href="{{.Path}}">{{.Name}}</a>
								{{else if hasSuffix .Name ".ogg"}}
									<audio controls preload="none">
										<source src="{{.Path}}" type="audio/ogg">
										Your browser does not support the audio tag.
									</audio><br>
									<a href="{{.Path}}">{{.Name}}</a>
								{{else if hasSuffix .Name ".opus"}}
									<audio controls preload="none">
										<source src="{{.Path}}" type="audio/opus">
										Your browser does not support the audio tag.
									</audio><br>
									<a href="{{.Path}}">{{.Name}}</a>
								{{else if or (hasSuffix .Name ".mp3") (hasSuffix .Name ".flac")}}
									<audio controls preload="none">
										<source src="{{.Path}}" type="audio/mpeg">
										<source src="{{.Path}}" type="audio/flac">
										Your browser does not support the audio tag.
									</audio><br>
									<a href="{{.Path}}">{{.Name}}</a>
								{{else if or (hasSuffix .Name ".txt") (hasSuffix .Name ".log")}}
									<a href="{{.Path}}" download>{{.Name}}</a><br>
									<textarea readonly rows="20" cols="80">{{ readFileContent .Path }}</textarea>
								{{else}}
									<a href="{{.Path}}">{{.Name}}</a>
								{{end}}
							{{else}}
								<a href="{{.Path}}">{{.Name}}</a>
							{{end}}
						{{end}}
					</td>
					<td>{{.Size}}</td>
				</tr>
				{{end}}
			</table>
			<p>Neekaru</p>
		</body>
		</html>
	`))

	err = tmpl.Execute(w, filesList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/", handler)

	log.Println("Server is running")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
