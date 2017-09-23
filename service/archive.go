package service

import(
  // "fmt"
  "io"
  "archive/zip"
  "net/http"
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
    for _, file := range a.Files {

      // println("Expprt zip archive: " + file.Url)

      f, err := z.Create(file.Url)
      if err != nil {
        return CommandError(http.StatusInternalServerError, err.Error())
      }

      // TODO: support full and public types
      if archive.Type == ArchiveSource {
        _, err = f.Write([]byte(file.Source(true)))
        if err != nil {
          return CommandError(http.StatusInternalServerError, err.Error())
        }
      }
    }

    err := z.Close()
    if err != nil {
      return CommandError(http.StatusInternalServerError, err.Error())
    }
  }
  return nil
}
