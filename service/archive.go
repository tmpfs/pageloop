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
  ArchiveSource = iota
  ArchivePublic
)

type ArchiveService struct {
  // Reference to the host
  Host *Host
}

type ArchiveRequest struct {
  // Name of the output file
  Name string
  // Reference to the target application
  Application *Application
  // Output stream
  Writer io.Writer
  // Type of archive to create. Full, source only or public only.
  Type int
}

// Export a zip archive of application files.
func (s *ArchiveService) Export(archive *ArchiveRequest, reply *ServiceReply) *StatusError {
  if _, a, err := LookupApplication(s.Host, archive.Application); err != nil {
    return err
  } else {
    z := zip.NewWriter(archive.Writer)

    // Send using in-memory file data
    if archive.Type == ArchiveSource {
      for _, file := range a.Files {
        f, err := z.Create(file.Url)
        if err != nil {
          return CommandError(http.StatusInternalServerError, err.Error())
        }

        _, err = f.Write([]byte(file.Source(true)))
        if err != nil {
          return CommandError(http.StatusInternalServerError, err.Error())
        }
      }

    // Walk the public files.
    } else {
      dir := a.PublicDirectory()
      err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
          return err
        }

        if path == dir || info.IsDir() {
          return nil
        }

        // Assuming POSIX style fs
        url := strings.TrimPrefix(path, dir)

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
    }

    err := z.Close()
    if err != nil {
      return CommandError(http.StatusInternalServerError, err.Error())
    }
  }
  return nil
}
