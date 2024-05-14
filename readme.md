## Setting up a Systemd Service for a Go Executable in Ubuntu

1. **Create the systemd service file:**
   ```bash
   sudo nano /etc/systemd/system/nik.service
   ```
2. **Add the following content to the service file, replacing /home/Lovely/Go-github-auto-update/nik with the actual path to your Go executable file:**

   ```bash
   [Unit]
   Description=Go Github Auto Update Service
   After=network.target

   [Service]
   Type=simple
   User=jsguru
   WorkingDirectory=/home/jsguru/Documents/Lovely/Go-github-auto-update
   ExecStart=/home/jsguru/Documents/Lovely/Go-github-auto-update/nik
   StandardOutput=append:/var/log/nik.log
   StandardError=inherit

   [Install]
   WantedBy=multi-user.target


   sudo systemctl daemon-reload
   ```

3. **Save the file and exit the text editor**
4. **Reload systemd to apply the changes:**
   ```bash
   sudo systemctl daemon-reload
   ```
5. **Start the service:**
   ```bash
   sudo systemctl start nik.service
   ```
6. **Check the status of the service:**
   ```bash
   sudo systemctl status nik.service
   ```
7. **Enable the service to start automatically on system boot:**
   ```bash
   sudo systemctl enable nik.service
   ```
8. **Check log:**
   ```bash
   cat /var/log/nik.log
   ```
