import { NavLink, Outlet } from "react-router-dom";

export function AppShell() {
  return (
    <div className="app-shell">
      <header className="topbar">
        <NavLink to="/" className="brand" end>
          <span className="brand-mark">
            <span>Forge</span>
          </span>
          <span className="brand-sub">Platform</span>
        </NavLink>
        <nav className="nav" aria-label="Principal">
          <NavLink to="/" end>
            Catalog
          </NavLink>
          <NavLink to="/request">Demande</NavLink>
          <NavLink to="/status">Suivi</NavLink>
          <NavLink to="/settings">Settings</NavLink>
        </nav>
      </header>
      <main className="main">
        <Outlet />
      </main>
    </div>
  );
}
