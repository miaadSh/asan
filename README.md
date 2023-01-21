# asan document

- how use asan cli

```bash
asan command options
```

---

- install command

```bash
asan install ProjectType ProjectName version env
```

- ProjectType: go - rust - laravel
- version: specfic version - default latest version
- env: docker / real mode environment

---

- check command

```bash
asan check ServiceName
```

- ServiceName: nginx, apache, composer, php, go, rust, mysql accourding os

---

- service-install command

```bash
asan service-install ServiceName
```

- ServiceName: nginx, apache, composer, php, go, rust, mysql accourding os

---

- service-upgrade command

```bash
asan service-upgrade ServiceName
```

- ServiceName: nginx, apache, composer, php, go, rust, mysql accourding os

---

- deploy command

```bash
asan deploy ProjectName
```
