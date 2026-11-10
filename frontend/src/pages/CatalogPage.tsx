import { useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { getPort, type CatalogItem } from "../lib/provisioning";

export function CatalogPage() {
  const [items, setItems] = useState<CatalogItem[]>([]);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  useEffect(() => {
    void getPort()
      .listCatalog()
      .then(setItems)
      .catch((err: unknown) =>
        setError(err instanceof Error ? err.message : String(err)),
      );
  }, []);

  return (
    <>
      <section className="hero">
        <h1>
          <em>Forge</em>
        </h1>
        <p>
          Portail self-service Platform Engineering — cataloguez un service,
          demandez une ressource, suivez son provisionnement.
        </p>
      </section>

      <section className="panel" aria-labelledby="catalog-title">
        <h2 id="catalog-title">Catalogue de services</h2>
        {error ? <p className="error">{error}</p> : null}
        <div className="catalog-grid">
          {items.map((item) => (
            <button
              key={item.id}
              type="button"
              className="catalog-item"
              disabled={!item.available}
              onClick={() => {
                if (item.id === "virtual-machine") navigate("/request");
              }}
            >
              <span className={`pill ${item.available ? "" : "soon"}`.trim()}>
                {item.available ? "Disponible" : "Bientôt"}
              </span>
              <h3>{item.title}</h3>
              <p>{item.description}</p>
              {item.available ? (
                <span className="mono" style={{ fontSize: "0.8rem" }}>
                  Demander →
                </span>
              ) : null}
            </button>
          ))}
        </div>
      </section>

      <p className="empty" style={{ marginTop: "1.25rem" }}>
        Déjà une demande ? Voir le{" "}
        <Link to="/status">suivi</Link>.
      </p>
    </>
  );
}
