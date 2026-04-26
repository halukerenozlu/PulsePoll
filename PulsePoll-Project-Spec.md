# PulsePoll-project-spec

## Amaç

Bu dosya, PulsePoll projesinin Türkçe ana referans özeti olarak tutulur.

Amaç:

- projeye ara verdikten sonra hızlı geri dönüş sağlamak
- ürün mantığını, teknik yönü ve çalışma modelini tek yerde toplamak
- agent rollerini ve proje akışını netleştirmek
- mevcut İngilizce contract dokümanlarının üst seviye özetini sunmak

Bu dosya üst seviye rehberdir.
Detay ve bağlayıcı teknik kontratlar hâlâ `docs/` altındaki dosyalardadır.

---

## Proje Özeti

PulsePoll, geçici yapıda çalışan bir anket platformu MVP'sidir.

Temel fikir:

- anketler kısa ömürlüdür
- oylama süresi sınırlıdır
- sonuç görünürlüğü ürün fazına bağlıdır
- sistem küçük, sade ve doğrulanabilir tutulur

Bu proje şu an private olarak geliştirilmektedir.

---

## Ana İlke

Bu projede kaynak gerçeklik dokümanlardır.

Eğer kod, plan, prompt, review notu veya başka bir çıktı aşağıdaki dosyalarla çelişirse dokümanlar kazanır:

- `docs/SPEC.md`
- `docs/API.md`
- `docs/DB.md`
- `docs/REDIS.md`
- `docs/verification.md`
- `docs/VERSION_PLAN.md`
- `docs/ROADMAP.md`
- `CHANGELOG.md`

Davranış değişecekse önce ilgili doküman güncellenmelidir.

---

## Teknoloji Yığını

### Backend

- Go
- Fiber

### Veritabanı

- PostgreSQL

### Geçici durum / oran sınırlama / kısa ömürlü kayıtlar

- Redis

### Frontend

- Next.js

### Yerel geliştirme

- Docker Compose

---

## Projenin Şu Anki Durumu

Aktif planlama modeli artık Version Milestone (Sürüm Kilometre Taşı) modelidir.

Mevcut tamamlanmış baseline:

- `v0.1.0` - Backend Foundation and Verification Baseline

Sıradaki planlama alanı repo durumuna göre:

- `v0.1.x` - küçük stabilizasyon, doküman temizliği, local doğrulama notları, API test rehberi cilası
- `v0.2.0` - kalan backend MVP endpoint/flow işlerinin tamamlanması

Önemli not:

- `v0.1.0`, eski planlama modelindeki tamamlanmış baseline işleri birleştirir.
- Bu, eski daha geniş backend feature/readiness kapsamının tamamen bittiği anlamına gelmez.
- Kalan backend ve API readiness işleri ilgili gelecek Version Milestone'lara taşınır.

---

## Version Milestone Modeli

Planlama terimleri:

- Version Milestone (Sürüm Kilometre Taşı): `v0.1.0`, `v0.2.0` gibi sürüm düzeyinde teslim hedefi.
- Work Item (İş Başlığı): bir Version Milestone içindeki anlamlı backend/frontend/docs/product hedefi.
- Implementation Slice (Uygulama Dilimi): bir Work Item altında Codex'in uygulayabileceği küçük, net kapsamlı iş parçası.

Kod işi başlamadan önce Version Milestone, Work Item ve Implementation Slice açık olmalıdır.
Codex aktif Implementation Slice kapsamını sessizce genişletmemelidir.

Aktif planlama kaynağı:

- `docs/VERSION_PLAN.md`

Yüksek seviye sıra:

- `v0.1.0` - Backend Foundation and Verification Baseline
- `v0.1.x` - Stabilization and Docs Cleanup
- `v0.2.0` - Backend Feature Completion
- `v0.3.0` - API Contract Readiness
- `v0.4.0` - Frontend Integration
- `v0.5.0` - End-to-End MVP Hardening

---

## Ürün Kuralları Özeti

### Anket fazları

Anketler üç temel durumda değerlendirilir:

- `VOTING`
- `RESULTS`
- `EXPIRED`

Bu terimler ürün/domain terminolojisidir ve proje planlama modelindeki eski faz terimleriyle karıştırılmamalıdır.

Genel mantık:

- `now < vote_ends_at` -> `VOTING`
- `vote_ends_at <= now < results_ends_at` -> `RESULTS`
- `now >= results_ends_at` -> `EXPIRED`

### Varsayılan süreler

MVP varsayılanları:

- `vote_ends_at = created_at + 24h`
- `results_ends_at = created_at + 48h`
- `retention_ends_at = created_at + 48h`

### Görünürlük türleri

- `PUBLIC`
- `UNLISTED`
- `PRIVATE_PIN`

### Sonuç görünürlüğü

- `OPEN_LIVE`
- `CLOSED_HIDDEN_UNTIL_END`

### Kullanıcı kuralları

- kayıtlı kullanıcılar anket oluşturabilir
- guest kullanıcılar anket oluşturamaz
- guest kullanıcılar oy vermek için gerekli çerezi kabul etmelidir
- guest kullanıcılar çerez kabul etmeden de gezinme ve izin verilen sonuçları görme hakkına sahiptir

### Oy kuralları

- `max_votes_per_user >= 1`
- `allow_vote_change_once` yalnızca `max_votes_per_user == 1` ise anlamlıdır
- oy değiştirme yalnızca `VOTING` fazında ve en fazla bir kez yapılabilir
- aynı kullanıcı veya guest, kural izin veriyorsa aynı seçeneğe birden fazla kez oy verebilir

### Moderasyon

MVP düzeyinde:

- anket oluşturulurken temel keyword filtresi vardır
- uygunsuz terim tespit edilirse oluşturma reddedilir
- report endpoint vardır

---

## Consent (Guest Voting) Mantığı

Guest kullanıcılar her zaman şunları yapabilir:

- public feed'i gezebilir
- anket detaylarını görebilir
- izin verilen sonuçları görebilir

Ama oy vermek veya oy değiştirmek için:

- gerekli servis çerezini kabul etmeleri gerekir
- bu çerez içinde kısa ömürlü bir `guest_id` tutulur

Bu yapı şu amaçlara hizmet eder:

- spam azaltma
- tekrar oy verme suistimalini sınırlama
- oy limiti uygulama
- bir kez oy değiştirme kuralını izleme
- kısa ömürlü PIN doğrulama durumunu hatırlama

---

## Veri Saklama Mantığı

### PostgreSQL

Kalıcı ve çekirdek veri burada tutulur.

Örnek tablolar:

- `users`
- `auth_sessions`
- `surveys`
- `survey_options`
- `reports`
- `feedback` (opsiyonel)

MVP yaklaşımı:

- ham oy event log tutulmaz
- aggregate sayılar tutulur
- `survey_options.vote_count` kritik sayaç alanıdır

### Redis

Geçici ve TTL (Time To Live - Yaşam Süresi) odaklı veriler burada tutulur.

Örnek kullanım alanları:

- vote receipts
- guest bazlı oy limiti takibi
- bir kez oy değiştirme durumu
- PIN doğrulama durumu
- brute-force önleme sayaçları
- rate limiting

---

## API Özeti

Base path:

- `/api/v1`

Önemli endpoint grupları:

- Auth
- Consent
- Surveys
- Feed
- PIN verify
- Vote / vote change
- Results
- Report
- Feedback (opsiyonel)

Önemli hata sınıfları:

- `400 BAD_REQUEST`
- `401 UNAUTHORIZED`
- `403 FORBIDDEN`
- `404 NOT_FOUND`
- `429 TOO_MANY_REQUESTS`
- `500 INTERNAL_SERVER_ERROR`

Özellikle takip edilen bazı hata kodları:

- `CONSENT_REQUIRED`
- `PIN_REQUIRED`
- `PHASE_NOT_VOTING`
- `VOTE_CHANGE_NOT_ALLOWED`

---

## Verification Yaklaşımı

Bu projede backend doğruluğu frontend'e bırakılmaz.

`docs/verification.md` şu amaçla vardır:

- backend'i UI olmadan doğrulamak
- local startup kontrolü yapmak
- `/health` durumunu doğrulamak
- endpoint success / failure senaryolarını kontrol etmek
- persistence etkilerini doğrulamak
- Version Milestone / Work Item / Implementation Slice bazında tekrar üretilebilir bir kontrol hattı oluşturmak

Ana fikir:

- frontend, backend'in ilk test edildiği yer olmamalıdır

---

## Agent Rolleri

### İnsan

- yönü onaylar
- commit ve tag sınırını belirler
- son kararı verir
- gerekirse local doğrulamayı kendisi de çalıştırır

### ChatGPT

- planlama yapar
- scope netleştirir
- Version Milestone / Work Item / Implementation Slice kurgusu çıkarır
- prompt üretir
- sıradaki adımı sadeleştirir

### Codex

Varsayılan implementer'dır.

Görevleri:

- onaylı Implementation Slice'ı uygular
- davranış değiştiyse test ekler veya günceller
- ilgili test/build komutlarını çalıştırır
- ne yaptığını ve neyi doğruladığını raporlar

### Gemini

İlk review ve özellikle frontend / product flow tarafında güçlü yardımcıdır.

Görevleri:

- Codex çıktısına first-pass review yapmak
- frontend yapısını değerlendirmek
- UX akışını eleştirmek
- gereksiz karmaşıklığı işaretlemek
- maintainability açısından not düşmek
- backend field, endpoint veya ürün davranışı uydurmadığını kontrol etmek

### Claude

Varsayılan implementer değildir.
Seçici ve yüksek değerli derin review aracı olarak kullanılır.

Daha çok şu alanlarda devreye girer:

- auth/session
- migration
- DB hassasiyeti
- security kritik kod
- karmaşık backend refactor
- önemli commit öncesi yüksek riskli inceleme

---

## Varsayılan Çalışma Akışı

Normal akış:

1. İnsan + ChatGPT Version Milestone, Work Item ve Implementation Slice belirler
2. Codex kodu uygular
3. Codex gerekli testleri ekler/günceller
4. Codex ilgili test/build adımlarını çalıştırır
5. Gemini first-pass review yapar
6. Codex gerekli düzeltmeleri yapar ve tekrar doğrular
7. İnsan sonucu inceler
8. ChatGPT commit/tag/sıradaki adım konusunda destek verir

Yüksek riskli işlerde ek adım:

- Gemini sonrası Claude selective deep review yapar
- gerekiyorsa Codex final patch uygular

---

## Test Politikası

Her değişiklik yeni test dosyası gerektirmez.
Ama her değişiklik uygun doğrulama gerektirir.

### Test eklenmesi gereken durumlar

- business logic değişikliği
- validation değişikliği
- error behavior değişikliği
- auth/session akışı
- route/handler davranışı
- bug fix
- kritik helper fonksiyonlar
- persistence kuralları

### Yeni test şart olmayabilecek durumlar

- docs-only değişiklik
- comment-only değişiklik
- davranış etkilemeyen küçük rename
- zaten mevcut testlerin kapsadığı mekanik refactor

Kural:

- davranış değiştiyse test ekle veya neden eklenmediğini açıkça yaz
- teslimden önce ilgili test/build komutlarını çalıştır

---

## Quality Bar

Minimum kabul seviyesi:

- açık isimlendirme
- net hata yönetimi
- deterministik davranış
- gereksiz scope genişlemesi olmaması
- migration gerekiyorsa commit edilmiş olması
- verification adımlarının tekrar üretilebilir olması
- ilgili test/build komutlarının çalıştırılmış olması

---

## Doküman Haritası

### Ürün ve davranış

- `docs/SPEC.md`

### API kontratı

- `docs/API.md`

### Veritabanı kontratı

- `docs/DB.md`

### Redis kontratı

- `docs/REDIS.md`

### Doğrulama akışı

- `docs/verification.md`
- `docs/API_TESTING.md`

### Planlama ve tarihçe

- `docs/VERSION_PLAN.md`
- `docs/ROADMAP.md`
- `CHANGELOG.md`

### Üst seviye özetler

- `docs/VISION.md`
- `docs/ARCHITECTURE.md`

---

## Local Çalıştırma Özeti

Servisleri başlatmak için örnek komut:

```bash
docker compose -p pulsepoll up --build
```

Health kontrolü:

```text
GET http://localhost:8080/health
```

Beklenen sağlıklı yanıt yapısı:

```json
{
  "db": "up",
  "ok": true,
  "redis": "up"
}
```

---

## Bu Dosya Nasıl Kullanılmalı

Bu dosya:

- hızlı geri dönüş rehberi
- proje hafıza dosyası
- üst seviye özel referans
- agent ve workflow özeti

olarak kullanılmalı.

Detay değişikliklerinde esas kaynak yine teknik kontrat dosyaları olmalıdır.
