package service

import(
  //"fmt"
  "os"
  "io"
  "io/ioutil"
  "strings"
  "archive/zip"
  "net/http"
  "path/filepath"
  . "github.com/tmpfs/pageloop/model"
  . "github.com/tmpfs/pageloop/util"
)

const(
  ArchiveFull = iota
  ArchiveSource
  ArchivePublic
)

type ArchiveService struct {
  // Reference to the host
  Host *Host
}

type ArchiveRequest struct {
  // Name of the output file
  Name string `json:"name"`
  // Type of archive to create. Full, source only or public only.
  Type int `json:"type"`
  // A reference to an application in the form: file://pageloop.com/{container}/{application}
  Ref string `json:"ref,omitempty"`
  // Output stream
  Writer io.Writer
}

// Export a zip archive of application files.
func (s *ArchiveService) Export(archive *ArchiveRequest, reply *ServiceReply) *StatusError {

  //fmt.Printf("%#v\n", archive)

  if archive.Writer == nil {
    return CommandError(http.StatusInternalServerError, "No writer for archive stream.")
  }

  // Add source files to the archive
  source := func (z *zip.Writer, app *Application, prefix string) *StatusError {
    for _, file := range app.Files {
      url := prefix + file.Url
      f, err := z.Create(url)
      if err != nil {
        return CommandError(http.StatusInternalServerError, err.Error())
      }

      // Send using in-memory file data
      _, err = f.Write([]byte(file.Source(true)))
      if err != nil {
        return CommandError(http.StatusInternalServerError, err.Error())
      }
    }
    return nil
  }

  // Add public files to the archive
  public := func (z *zip.Writer, app *Application, prefix string) *StatusError {
    // Walk the public files.
    dir := app.PublicDirectory()
    err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
      if err != nil {
        return err
      }

      if path == dir || info.IsDir() {
        return nil
      }

      // Assuming POSIX style fs
      url := strings.TrimPrefix(path, dir)
      url = prefix + url

      f, err := z.Create(url)
      if err != nil {
        return CommandError(http.StatusInternalServerError, err.Error())
      }

      // TODO: stream content from disc
      if content, err := ioutil.ReadFile(path); err != nil {
        return err
      } else {
        _, err = f.Write(content)
        if err != nil {
          return CommandError(http.StatusInternalServerError, err.Error())
        }
      }

      return nil
    })
    if err != nil {
      return CommandError(http.StatusInternalServerError, err.Error())
    }
    return nil
  }

  ref := &AssetReference{}
  ref.ParseUrl(archive.Ref)
  if _, a, err := ref.FindApplication(s.Host); err != nil {
    return err
  } else {
    z := zip.NewWriter(archive.Writer)
    if archive.Type == ArchiveFull {
      if err := source(z, a, "/source"); err != nil {
        return err
      }
      if err := public(z, a, "/public"); err != nil {
        return err
      }
    } else if archive.Type == ArchiveSource {
      if err := source(z, a, ""); err != nil {
        return err
      }
    } else {
      if err := public(z, a, ""); err != nil {
        return err
      }
    }

    err := z.Close()
    if err != nil {
      return CommandError(http.StatusInternalServerError, err.Error())
    }
  }
  return nil
}
