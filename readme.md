## Automatic Github repository Update\

1. **Main Functionalities**

   - Monitor updates on the GitHub repository by comparing commit hashes.
   - Automatically pull or clone (for the first time) the repository when updates are monitored.
   - Execute a running file in it.
   - Update the local commit hash value to _LastSuccessCommitHash.txt_ file when execution succeeds.
   - Revert to the last successful commit when execution fails.
   - Remember the last failed commit hash to _LastFailedCommitHash.txt_ file for skipping in future attempts.
   - Provide a REST API for getting and setting _supperFlag_ which turns the thread on/off.

2. **Update the following fields with your info in main_public.go file:**

   ```bash
   const (
      owner = "LovelyDev829"
      repo = "nik-hello-world"
      interval = 5 * time.Second
      localPath = "../repo-pulled/"
      execCommand = localPath + "hello-world.exe"
      successHashFile = "LastSuccessCommitHash.txt"
      failedHashFile = "LastFailedCommitHash.txt"
   )
   ```

3. **Run and Test:**
   ```bash
   go main_public.go
   ```

## Setting up a Systemd Service for a Go Executable in Ubuntu

1. **Build the executable file first:**

   ```bash
   go build main_public.go
   ```

1. **Create the systemd service file:**
   ```bash
   sudo nano /etc/systemd/system/nik.service
   ```
1. **Add the following content to the service file, replacing /home/Lovely/Go-github-auto-update/nik with the actual path to your Go executable file:**

   _(This is an example and you can update WorkingDirectory and ExecStat below)_

   ```bash
   [Unit]
   Description=Go Github Auto Update Service
   After=network.target

   [Service]
   Type=simple
   User=jsguru
   WorkingDirectory=/home/jsguru/Documents/Lovely/Go-github-auto-update
   ExecStart=/home/jsguru/Documents/Lovely/Go-github-auto-update/main_public
   StandardOutput=append:/var/log/nik.log
   StandardError=inherit

   [Install]
   WantedBy=multi-user.target


   sudo systemctl daemon-reload
   ```

1. **Save the file and exit the text editor**
1. **Reload systemd to apply the changes:**
   ```bash
   sudo systemctl daemon-reload
   ```
1. **Start the service:**
   ```bash
   sudo systemctl start nik.service
   ```
1. **Check the status of the service:**
   ```bash
   sudo systemctl status nik.service
   ```
1. **Enable the service to start automatically on system boot:**
   ```bash
   sudo systemctl enable nik.service
   ```
1. **Check REST API:**
   ```bash
   [GET]  localhost:8080/get-supper-flag
   [POST] localhost:8080/update-supper-flag with {"supperFlag": true} as BODY
   ```
