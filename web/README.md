# QuickRT — Alpine.js Edition

React → Alpine.js conversion. No Vite, no React runtime.  
Built to embed directly into a **Gin (Go)** backend.

---

## Stack

| Layer | Tool | Replaces |
|---|---|---|
| Reactivity | Alpine.js v3 | React + useState/useEffect |
| Global state | `Alpine.store()` | React Context / hooks |
| Persistence | `@alpinejs/persist` | localStorage + useEffect sync |
| DOM diffing | `@alpinejs/morph` | Virtual DOM |
| Bundler | esbuild | Vite |
| Styles | Tailwind CSS v4 | same |
| Types | TypeScript (type-check only) | same |
| HTTP | native `fetch` | axios |
| Templates | Gin HTML templates | JSX |

---

## Project Structure

```
quickrt-alpine/
├── templates/                  ← Gin serves these
│   ├── auth.html               ← /auth page
│   └── dashboard.html          ← /user page
├── static/
│   ├── css/
│   │   ├── app.css             ← Tailwind source
│   │   └── out.css             ← Compiled (gitignore this)
│   └── js/
│       ├── app.ts              ← Entry point — registers everything
│       ├── app.js              ← Compiled bundle (gitignore this)
│       ├── stores/
│       │   ├── auth.store.ts   ← AuthCtx + token logic
│       │   ├── ui.store.ts     ← Toasts + theme
│       │   └── dashboard.store.ts ← DashboardLayoutICtx + viewNavTree
│       ├── components/
│       │   ├── auth/
│       │   │   ├── auth-flow.component.ts      ← pages/Auth.tsx
│       │   │   ├── auth-entry.component.ts     ← AuthEntry.tsx
│       │   │   ├── auth-otp.component.ts       ← AuthOTP.tsx
│       │   │   ├── auth-profile.component.ts   ← AuthProfileSetup.tsx
│       │   │   └── auth-success.component.ts   ← AuthSuccess.tsx
│       │   └── dashboard/
│       │       ├── dashboard-layout.component.ts    ← DashboardLayoutI.tsx
│       │       ├── retractable-sidebar.component.ts ← RetractableSidebar.tsx
│       │       ├── mobile-nav.component.ts          ← MobileNav.tsx
│       │       └── top-nav-link.component.ts        ← TopNavLink.tsx
│       └── lib/
│           ├── api.ts          ← axiosInstance.ts (native fetch)
│           ├── helpers.ts      ← lib/helpers.ts (no React deps)
│           └── app-config.ts   ← lib/appConfig.ts
├── gin_controllers_example.go  ← How to wire pages in Gin
├── package.json
└── tsconfig.json
```

---

## Getting Started

### 1. Install dependencies
```bash
npm install
```

### 2. Build for production
```bash
npm run build
# Outputs: static/js/app.js + static/css/out.css
```

### 3. Dev mode (watch)
```bash
npm run dev
# Watches TS and CSS simultaneously
```

### 4. Wire into Gin (see gin_controllers_example.go)
```go
r.Static("/static", "./static")
r.LoadHTMLGlob("templates/*")
r.GET("/auth", controllers.ServeAuth)
r.GET("/user", controllers.ServeUserDashboard)
```

> Update your HTML templates to use `out.css` instead of `app.css` once you run a build.

---

## React → Alpine Mapping

| React concept | Alpine equivalent |
|---|---|
| `useState` | `x-data` local property |
| `useEffect` | `init()` + `destroy()` in `Alpine.data()` |
| `useContext` | `Alpine.store()` / `$store.name` |
| `React.createContext` | `Alpine.store('name', {...})` |
| Custom hook | Function in `lib/` called inside `Alpine.data()` |
| Component props | Args to `Alpine.data('name', (prop) => ({...}))` |
| `AnimatePresence` | `x-transition` directives |
| `framer-motion` | CSS transitions + `x-transition` |
| `react-router-dom` | Gin routes serve full pages |

---

## Adding a New Page / View

1. Create `templates/my-page.html`
2. Add a Gin route: `r.GET("/my-page", controllers.ServeMyPage)`
3. If the page needs a new Alpine component, add `static/js/components/my-page.component.ts` and import it in `app.ts`
4. Run `npm run build`

---

## Adding a New Dashboard View

In `templates/dashboard.html`, add:
```html
<div x-show="activeView === 'MyView'" x-transition>
  <!-- your content -->
</div>
```

In `gin_controllers_example.go`, add to `dashboardNavLinks`:
```go
{Path: "/my-view", Label: "MyView", Icon: iconMyView},
```
