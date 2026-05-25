# PulsePoll — Docker Kısa Rehber (unutmamak için)

Bu notlar `docker compose -p pulsepoll ...` kullanımı içindir.

---

## 1) Durum kontrol
```bash
docker compose -p pulsepoll ps
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
```

Health testi:
```bash
curl http://localhost:8080/health
```
Tarayıcı:
- http://localhost:8080/health

---

## 2) Log ekranından çıkma (Detach)
Compose “Attaching to ...” modundayken:
- `d` → Detach (servisler çalışmaya devam eder)
- `Ctrl + C` → compose’u durdurur (container’lar kapanır)

---

## 3) Stop vs Down (farkı)
### `stop`
- Container’ları **durdurur**, container objeleri yerinde kalır.
- Hızlı geri dönüş için uygundur.
```bash
docker compose -p pulsepoll stop
docker compose -p pulsepoll start
```

### `down`
- Container + network **kaldırılır/silinir**.
- **Volume’lar kalır** (DB verisi genelde durur) — *`-v` kullanmadıkça*.
```bash
docker compose -p pulsepoll down
```

---

## 4) Senin senaryon (diğer projelerde de Postgres var)
Port çakışmalarını (özellikle 5432) yaşamamak için:
- **Başka projeye geçerken:** `down`
- **Geri dönerken:** `up -d`

### Başka projeye geç
```bash
docker compose -p pulsepoll down
```

### PulsePoll’a geri dön
```bash
docker compose -p pulsepoll up -d
```

Değişiklik yaptıysan:
```bash
docker compose -p pulsepoll up -d --build
```

---

## 5) Tam sıfırlama (DB dahil her şey gider)
Volume’ları da siler (Postgres/Redis verisi silinir):
```bash
docker compose -p pulsepoll down -v
```

---

## 6) Log izleme
Arka planda çalışırken logları takip et:
```bash
docker compose -p pulsepoll logs -f
```

---

## 7) “5432’yi kim kullanıyor?” hızlı kontrol
```bash
docker ps --format "table {{.Names}}\t{{.Ports}}"
```

---

## 8) Genel temizlik (isteğe bağlı)
Kullanılmayan image/container/network temizler:
```bash
docker system prune
```
Not: Sonra bazı image’ları yeniden indirmen gerekebilir.
