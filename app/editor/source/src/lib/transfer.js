/**
 *  Helper class for performing chunked file transfers.
 *
 *  By chunk we mean not attempting to upload all files at
 *  once nor only one at a time.
 *
 *  When a file list is received it is broken into chunks
 *  based on the concurrentTransfers variable.
 *
 *  Each chunk is then processed in series whilst the chunk
 *  files are uploaded concurrently.
 */
class Transfer {
  constructor (client) {
    this.client = client

    // List of all upload transfers
    this.transfers = []

    // Number of files in each transfer chunk
    this.concurrentTransfers = 3

    // The list of files currently being transferred
    this.currentTransfer = []
  }

  upload () {
    if (this.transfers.length) {
      let amount = Math.floor(this.transfers.length / this.concurrentTransfers)
      if (this.transfers.length % this.concurrentTransfers !== 0) {
        amount++
      }

      let chunks = []
      let i, ind, len
      for (i = 0; i < amount; i++) {
        ind = i * this.concurrentTransfers
        len = Math.min(this.transfers.length, ind + this.concurrentTransfers)
        chunks.push(this.transfers.slice(ind, len))
      }

      // Transfer a single chunk
      const transfer = (chunk, done) => {
        return new Promise((resolve, reject) => {
          let loaded = 0
          chunk.forEach((file) => {
            this.client.upload(file).then((file) => {
              loaded++
              if (loaded === chunk.length) {
                // Process next chunk
                if (chunks.length) {
                  this.currentTransfer = chunks.shift()
                  resolve(transfer(this.currentTransfer, done))
                // All done, upload completed
                } else {
                  done(this.transfers)
                }
              }
            })
            .catch(reject)
          })
        })
      }
      this.currentTransfer = chunks.shift()
      return new Promise((resolve, reject) => {
        transfer(this.currentTransfer, (files) => {
          this.transfers = []
          this.currentTransfer = []
          resolve(files)
        })
        .catch((err) => {
          this.transfers = []
          this.currentTransfer = []
          reject(err)
        })
      })
    }
  }

}

export default Transfer
