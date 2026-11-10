import { useState } from "react";
import type { FormEvent } from "react";
import { getPort, type PortalSettings, type ProviderMode } from "../lib/provisioning";

function isMode(value: string): value is ProviderMode {
  return value === "mock" || value === "crossplane";
}

export function SettingsPage() {
  const [form, setForm] = useState<PortalSettings>(() => getPort().getSettings());
  const [saved, setSaved] = useState(false);

  function onSubmit(e: FormEvent) {
    e.preventDefault();
    getPort().saveSettings({
      apiBaseUrl: form.apiBaseUrl.trim() || "/api/v1",
      mode: form.mode,
    });
    setSaved(true);
    setTimeout(() => setSaved(false), 2500);
  }

  return (
    <>
      <section className="hero">
        <h1>
          <em>Settings</em>
        </h1>
        <p>
          Aucun secret n'est stocké par le portail. L'authentification auprès du
          cluster est portée par le BFF via son ServiceAccount.
        </p>
      </section>

      <form className="panel" onSubmit={onSubmit}>
        <h2>Provisionnement</h2>
        <div className="form-grid">
          <div className="field">
            <label htmlFor="mode">mode</label>
            <select
              id="mode"
              value={form.mode}
              onChange={(e) =>
                setForm({
                  ...form,
                  mode: isMode(e.target.value) ? e.target.value : "crossplane",
                })
              }
            >
              <option value="crossplane">crossplane (BFF → cluster)</option>
              <option value="mock">mock (local, sans cluster)</option>
            </select>
          </div>

          <div className="field">
            <label htmlFor="apiBaseUrl">API du BFF</label>
            <input
              id="apiBaseUrl"
              className="mono"
              placeholder="/api/v1"
              value={form.apiBaseUrl}
              onChange={(e) => setForm({ ...form, apiBaseUrl: e.target.value })}
            />
          </div>
        </div>

        {saved ? <p className="toast">Settings enregistrés.</p> : null}

        <div className="actions">
          <button className="btn btn-primary" type="submit">
            Enregistrer
          </button>
        </div>
      </form>
    </>
  );
}
