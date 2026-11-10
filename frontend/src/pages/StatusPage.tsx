import { useCallback, useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { StatusBadge } from "../components/StatusBadge";
import {
  getPort,
  type ProviderMode,
  type ProvisioningRequest,
} from "../lib/provisioning";

/** Crossplane réconcilie en continu : on rafraîchit le suivi périodiquement. */
const POLL_INTERVAL_MS = 5000;

export function StatusPage() {
  const [requests, setRequests] = useState<ProvisioningRequest[]>([]);
  const [mode, setMode] = useState<ProviderMode>("crossplane");
  const [error, setError] = useState<string | null>(null);

  const refresh = useCallback(async () => {
    const port = getPort();
    setMode(port.getSettings().mode);
    try {
      setRequests(await port.listRequests());
      setError(null);
    } catch (err) {
      setError(err instanceof Error ? err.message : String(err));
    }
  }, []);

  useEffect(() => {
    void refresh();
    const timer = setInterval(() => void refresh(), POLL_INTERVAL_MS);
    return () => clearInterval(timer);
  }, [refresh]);

  async function advance(id: string) {
    await getPort().advanceMockStatus(id);
    await refresh();
  }

  async function destroy(name: string) {
    if (!confirm(`Supprimer la VM « ${name} » et ses ressources Scaleway ?`)) {
      return;
    }
    try {
      await getPort().deleteRequest(name);
      await refresh();
    } catch (err) {
      setError(err instanceof Error ? err.message : String(err));
    }
  }

  return (
    <>
      <section className="hero">
        <h1>
          Suivi des <em>demandes</em>
        </h1>
        <p>
          {mode === "mock"
            ? "Statuts simulés localement. Mode actuel : "
            : "Statuts lus depuis les composites Crossplane. Mode actuel : "}
          <span className="mono">{mode}</span>.
        </p>
      </section>

      <section className="panel">
        <h2>Demandes</h2>
        {error ? <p className="error">{error}</p> : null}

        {requests.length === 0 ? (
          <p className="empty">
            Aucune demande. Commencez par le <Link to="/">catalogue</Link>.
          </p>
        ) : (
          <div className="request-list">
            {requests.map((r) => (
              <article key={r.id} className="request-row">
                <div>
                  <h3 className="mono">{r.parameters.vmName}</h3>
                  <p className="request-meta mono">
                    {r.parameters.machineType} · {r.parameters.zone} ·{" "}
                    {r.parameters.diskSizeGb} Go · {r.parameters.imageFamily}
                  </p>
                  {r.publicIp ? (
                    <p className="request-meta mono">IP {r.publicIp}</p>
                  ) : null}
                  {r.message ? <p className="request-meta">{r.message}</p> : null}
                  <p className="request-meta">
                    {new Date(r.createdAt).toLocaleString("fr-FR")}
                  </p>
                </div>

                <div className="row-actions">
                  <StatusBadge status={r.status} />
                  {mode === "mock" &&
                  r.status !== "ready" &&
                  r.status !== "failed" ? (
                    <button
                      type="button"
                      className="btn btn-ghost"
                      style={{ padding: "0.35rem 0.75rem", fontSize: "0.8rem" }}
                      onClick={() => void advance(r.id)}
                    >
                      Avancer statut
                    </button>
                  ) : null}
                  <button
                    type="button"
                    className="btn btn-ghost"
                    style={{ padding: "0.35rem 0.75rem", fontSize: "0.8rem" }}
                    onClick={() => void destroy(r.parameters.vmName)}
                  >
                    Supprimer
                  </button>
                </div>
              </article>
            ))}
          </div>
        )}
      </section>
    </>
  );
}
