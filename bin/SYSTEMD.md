# Systemd Configuration

## Files

```bash
root@epimethean:/etc/systemd/system# pwd
/etc/systemd/system

root@epimethean:/etc/systemd/system# ll epimethean*
-rw-r--r-- 1 root root 241 Oct  4 06:18 epimethean.service
-rw-r--r-- 1 root root 188 Oct 13 22:33 epimethean-tasks.timer
-rw-r--r-- 1 root root 421 Sep 30 22:06 epimethean-tasks.service

root@epimethean:/var/www/dev.epimethean/bin# pwd
/var/www/dev.epimethean/bin

root@epimethean:/var/www/dev.epimethean/bin# ll
total 24408
drwxrwxr-x 2 mdhender mdhender     4096 Oct 13 22:24 ./
drwxr-xr-x 6 mdhender mdhender     4096 Oct 13 21:57 ../
-rwxr-xr-x 1 mdhender mdhender  6957655 Oct 13 22:11 epimethean*
-rwxr-xr-x 1 mdhender mdhender      439 Oct  4 06:07 epimethean.service.sh*
-rwxr-xr-x 1 mdhender mdhender     5561 Oct 13 22:24 epimethean.sh*
```

## Install

```bash
root@epimethean:/etc/systemd/system# systemctl daemon-reload

root@epimethean:/etc/systemd/system# systemctl enable epimethean.service

root@epimethean:/etc/systemd/system# systemctl status epimethean.service
● epimethean.service - Epimethean Web server
     Loaded: loaded (/etc/systemd/system/epimethean.service; enabled; preset: enabled)
     Active: active (running) since Sun 2024-10-13 22:46:52 UTC; 4h 37min ago
   Main PID: 1778 (epimethean)
      Tasks: 7 (limit: 1112)
     Memory: 3.7M (peak: 6.2M)
        CPU: 153ms
     CGroup: /system.slice/epimethean.service
             └─1778 /var/www/dev.epimethean/bin/epimethean serve --database /home/epimethean/data/epimethean.db

root@epimethean:/etc/systemd/system# systemctl enable epimethean.service
○ epimethean.service - Run epimethean every 1 minutes
     Loaded: loaded (/etc/systemd/system/epimethean.service; static)
     Active: inactive (dead) since Mon 2024-10-14 03:25:28 UTC; 3s ago
   Duration: 57ms
TriggeredBy: ● epimethean.timer
    Process: 8709 ExecStart=/var/www/dev.epimethean/bin/epimethean.service.sh (code=exited, status=0/SUCCESS)
   Main PID: 8709 (code=exited, status=0/SUCCESS)
        CPU: 35ms

root@epimethean:/etc/systemd/system# systemctl enable epimethean.timer
● epimethean.timer - Run epimethean.service.sh 30 seconds after it finishes
     Loaded: loaded (/etc/systemd/system/epimethean.timer; enabled; preset: enabled)
     Active: active (waiting) since Sun 2024-10-13 22:33:51 UTC; 4h 52min ago
    Trigger: Mon 2024-10-14 03:25:58 UTC; 2s left
   Triggers: ● epimethean.service
```

## Monitor

```bash
root@epimethean:/etc/systemd/system# journalctl -f -u epimethean.service

root@epimethean:/etc/systemd/system# journalctl -f -u epimethean.service

root@epimethean:/etc/systemd/system# journalctl -f -u epimethean.timer
```

