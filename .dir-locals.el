((nil
  (compile-command . "CGO_ENABLED=0 go build -ldflags='-s -w' && gzip -c cmdsvc | ssh ansible.lucca.local 'sudo -u ansible /bin/sh -c \"gunzip - > /srv/ansible/cmdsvc && chmod +x /srv/ansible/cmdsvc\"'")))
