# اسکرام پوکر فارسی

یک MVP ساده، فارسی‌محور و production-minded برای Scrum Poker با:

- بک‌اند `Go + Fiber`
- فرانت‌اند `React + TypeScript + Vite + Tailwind`
- دیتابیس `MySQL`
- همگام‌سازی زنده با `Redis + WebSocket`

این پروژه از روز اول RTL، متن فارسی، تم تیره/روشن، رأی‌گیری مخفی تا رأی آخر، reveal خودکار، و تاریخچه هر تسک را پشتیبانی می‌کند.

## قابلیت‌های MVP

- ساخت اتاق و تبدیل سازنده به میزبان
- ورود سایر اعضا با کد اتاق و نام نمایشی
- فقط یک تسک فعال در هر اتاق
- رأی‌گیری با کارت‌های ثابت `0.5, 1, 2, 3, 5, 8, 13, 21, 34, ?, ∞, Coffee`
- پنهان ماندن رأی‌ها تا رأی دادن همه اعضای فعال
- reveal خودکار به‌محض کامل شدن رأی‌ها
- محاسبه میانگین فقط از رأی‌های عددی
- ذخیره کامل رأی‌ها و میانگین در تاریخچه
- پایداری سشن هر عضو با توکن opaque

## ساختار پروژه

```text
.
├── backend
│   ├── cmd/server
│   └── internal
│       ├── cache
│       ├── config
│       ├── http
│       ├── models
│       ├── service
│       ├── store
│       └── ws
├── frontend
│   └── src
│       ├── app
│       ├── components
│       ├── features
│       ├── hooks
│       ├── lib
│       └── styles
└── docker-compose.yml
```

## اجرای سریع با Docker Compose

```bash
docker compose up --build
```

سرویس‌ها:

- فرانت‌اند: [http://localhost:5173](http://localhost:5173)
- بک‌اند: [http://localhost:8080](http://localhost:8080)
- MySQL: `localhost:3306`
- Redis: `localhost:6379`

## اجرای محلی بدون Docker

### 1. اجرای MySQL و Redis

می‌توانید از Docker فقط برای زیرساخت استفاده کنید:

```bash
docker compose up mysql redis -d
```

### 2. بک‌اند

```bash
cd backend
cp .env.example .env
go run ./cmd/server
```

نکته: migrationها در زمان بالا آمدن سرور به‌صورت خودکار اجرا می‌شوند.

### 3. فرانت‌اند

```bash
cd frontend
cp .env.example .env
npm install
npm run dev
```

## تست و build

### بک‌اند

```bash
cd backend
go test ./...
go build ./...
```

### فرانت‌اند

```bash
cd frontend
npm test
npm run build
```

## مدل داده

### MySQL

داده‌های durable در MySQL نگه‌داری می‌شوند:

- `rooms`
- `members`
- `tasks`
- `votes`

### Redis

Redis فقط برای stateهای ephemeral و real-time استفاده می‌شود:

- presence اعضای اتاق
- active member tracking
- vote progress counter
- current task snapshot cache
- pub/sub برای fanout رویدادهای WebSocket

اگر Redis یا سرور restart شود، داده‌های نهایی کسب‌وکاری از MySQL قابل بازیابی هستند.

## قرارداد API

### REST

- `POST /api/rooms`
  بدنه: `{ "displayName": "مرصاد" }`
- `POST /api/rooms/join`
  بدنه: `{ "displayName": "ندا", "code": "ABC123" }`
- `GET /api/rooms/:code`
  هدر: `Authorization: Bearer <token>`
- `POST /api/rooms/:code/tasks`
  هدر: `Authorization: Bearer <token>`
  بدنه: `{ "title": "پیاده‌سازی صفحه گزارش" }`
- `POST /api/rooms/:code/tasks/:taskId/votes`
  هدر: `Authorization: Bearer <token>`
  بدنه: `{ "value": "8" }`
- `GET /api/rooms/:code/history`
  هدر: `Authorization: Bearer <token>`

### پاسخ ایجاد/ورود

سرور بعد از ایجاد یا ورود اتاق، این سشن را برمی‌گرداند:

```json
{
  "roomCode": "ABC123",
  "token": "opaqueRandomToken",
  "memberId": 1
}
```

## قرارداد WebSocket

Endpoint:

```text
ws://localhost:8080/ws/rooms/:code?token=<session-token>
```

کلاینت هر ۲۰ ثانیه پیام heartbeat می‌فرستد:

```json
{ "type": "heartbeat" }
```

رویدادهای سرور:

- `system.connected`
- `room.created`
- `member.joined`
- `member.left`
- `task.created`
- `vote.progress`
- `results.revealed`
- `history.updated`

کلاینت روی هر رویداد relevant، وضعیت اتاق را دوباره از REST می‌خواند تا state شخصی‌سازی‌شده هر عضو مثل `myVote` هم درست بماند.

## منطق reveal و میانگین

- فقط یک تسک active با وضعیت `voting` یا `revealed` در هر اتاق وجود دارد.
- فقط وقتی همه اعضای active رأی داده‌اند، reveal انجام می‌شود.
- رأی‌های `Coffee`، `∞` و `?` ذخیره می‌شوند اما در میانگین دخالت ندارند.
- اگر هیچ رأی عددی وجود نداشته باشد، خروجی میانگین: `بدون برآورد عددی`

## تجربه کاربری فارسی

- `lang="fa"` و `dir="rtl"` از ریشه HTML
- فونت اصلی: `Vazirmatn`
- نمایش room code و توکن‌های فنی با فونت mono
- تم تیره پیش‌فرض با ذخیره در `localStorage`
- تم روشن با همان سطح polish، نه نسخه فراموش‌شده

## فایل‌های مهم

- بک‌اند ورودی اصلی: `backend/cmd/server/main.go`
- منطق دومین: `backend/internal/service/room_service.go`
- منطق رأی و میانگین: `backend/internal/service/scoring.go`
- API و WebSocket: `backend/internal/http/router.go`
- صفحه اصلی فرانت: `frontend/src/features/home/home-page.tsx`
- صفحه اتاق: `frontend/src/features/room/room-page.tsx`
- وضعیت‌های رأی‌گیری: `frontend/src/features/room/components/task-stage.tsx`

## وضعیت راستی‌آزمایی

این دستورات روی همین workspace با موفقیت اجرا شدند:

- `cd backend && go test ./...`
- `cd backend && go build ./...`
- `cd frontend && npm test`
- `cd frontend && npm run build`
