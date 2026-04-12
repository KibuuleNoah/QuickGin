import Alpine from "alpinejs";
import persist from "@alpinejs/persist";
import morph from "@alpinejs/morph";

// Stores
import "./stores/auth.store";
import "./stores/ui.store";
import "./stores/dashboard.store";

// Components
import "./components/auth/auth-flow.component";
import "./components/auth/auth-entry.component";
import "./components/auth/auth-otp.component";
import "./components/auth/auth-profile.component";
import "./components/auth/auth-success.component";

// Dashboard
import "./components/dashboard/dashboard-layout.component";
import "./components/dashboard/retractable-sidebar.component";
import "./components/dashboard/mobile-nav.component";
import "./components/dashboard/top-nav-link.component";

Alpine.plugin(persist);
Alpine.plugin(morph);

// Make Alpine available globally for debugging
(window as any).Alpine = Alpine;

Alpine.start();
