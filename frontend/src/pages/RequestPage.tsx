import { useEffect, useState } from "react";
import type { FormEvent } from "react";
import { Link, useNavigate } from "react-router-dom";
import {
  getPort,
  ValidationError,
  type FormOptions,
  type VmParameters,
} from "../lib/provisioning";

const DEFAULTS: VmParameters = {
  vmName: "",
  machineType: "DEV1-S",
  zone: "fr-par-1",
  imageFamily: "ubuntu_jammy",
  diskSizeGb: 20,
  rootVolumeType: "sbs_volume",
  dataDiskSizeGb: 0,
  staticIp: false,
};

const FALLBACK_OPTIONS: FormOptions = {
  machineTypes: [DEFAULTS.machineType],
  zones: [DEFAULTS.zone],
  imageFamilies: [DEFAULTS.imageFamily],
  rootVolumeTypes: [DEFAULTS.rootVolumeType],
  minDiskSizeGb: 10,
  maxDiskSizeGb: 1000,
};

export function RequestPage() {
  const [form, setForm] = useState<VmParameters>(DEFAULTS);
  const [options, setOptions] = useState<FormOptions>(FALLBACK_OPTIONS);
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({});
  const [error, setError] = useState<string | null>(null);
  const [busy, setBusy] = useState(false);
  const [toast, setToast] = useState<string | null>(null);
  const navigate = useNavigate();

  // Les valeurs autorisées viennent de l'API : la XRD reste la seule source de vérité.
  useEffect(() => {
    void getPort()
      .getFormOptions()
      .then(setOptions)
      .catch((err: unknown) =>
        setError(err instanceof Error ? err.message : String(err)),
      );
  }, []);

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    setError(null);
    setFieldErrors({});
    setToast(null);

    const name = form.vmName.trim().toLowerCase();
    if (!/^[a-z]([a-z0-9-]{1,61}[a-z0-9])?$/.test(name)) {
      setFieldErrors({
        vmName:
          "minuscules, chiffres et tirets ; doit commencer par une lettre (2 à 63 caractères).",
      });
      return;
    }

    setBusy(true);
    try {
      await getPort().createRequest({ ...form, vmName: name });
      setToast("Demande créée. Redirection vers le suivi…");
      setTimeout(() => navigate("/status"), 700);
    } catch (err) {
      if (err instanceof ValidationError) {
        setFieldErrors(err.fields);
        setError(err.message);
      } else {
        setError(err instanceof Error ? err.message : String(err));
      }
    } finally {
      setBusy(false);
    }
  }

  return (
    <>
      <section className="hero">
        <h1>
          Demande <em>VM</em>
        </h1>
        <p>
          Champs alignés sur la XRD Crossplane{" "}
          <span className="mono">VirtualMachine</span>. La soumission crée un
          composite dans le cluster — voir <Link to="/settings">Settings</Link>.
        </p>
      </section>

      <form className="panel" onSubmit={onSubmit}>
        <h2>Paramètres</h2>
        <div className="form-grid">
          <div className="field">
            <label htmlFor="vmName">vmName</label>
            <input
              id="vmName"
              className="mono"
              required
              value={form.vmName}
              onChange={(e) => setForm({ ...form, vmName: e.target.value })}
              placeholder="my-dev-server"
            />
            {fieldErrors.vmName ? (
              <p className="error">{fieldErrors.vmName}</p>
            ) : null}
          </div>

          <div className="field">
            <label htmlFor="machineType">machineType</label>
            <select
              id="machineType"
              value={form.machineType}
              onChange={(e) => setForm({ ...form, machineType: e.target.value })}
            >
              {options.machineTypes.map((type) => (
                <option key={type} value={type}>
                  {type}
                </option>
              ))}
            </select>
            {fieldErrors.machineType ? (
              <p className="error">{fieldErrors.machineType}</p>
            ) : null}
          </div>

          <div className="field">
            <label htmlFor="zone">zone</label>
            <select
              id="zone"
              value={form.zone}
              onChange={(e) => setForm({ ...form, zone: e.target.value })}
            >
              {options.zones.map((zone) => (
                <option key={zone} value={zone}>
                  {zone}
                </option>
              ))}
            </select>
            {fieldErrors.zone ? <p className="error">{fieldErrors.zone}</p> : null}
          </div>

          <div className="field">
            <label htmlFor="imageFamily">imageFamily</label>
            <select
              id="imageFamily"
              value={form.imageFamily}
              onChange={(e) => setForm({ ...form, imageFamily: e.target.value })}
            >
              {options.imageFamilies.map((image) => (
                <option key={image} value={image}>
                  {image}
                </option>
              ))}
            </select>
            {fieldErrors.imageFamily ? (
              <p className="error">{fieldErrors.imageFamily}</p>
            ) : null}
          </div>

          <div className="field">
            <label htmlFor="diskSizeGb">diskSizeGb</label>
            <input
              id="diskSizeGb"
              type="number"
              min={options.minDiskSizeGb}
              max={options.maxDiskSizeGb}
              value={form.diskSizeGb}
              onChange={(e) =>
                setForm({ ...form, diskSizeGb: Number(e.target.value) })
              }
            />
            {fieldErrors.diskSizeGb ? (
              <p className="error">{fieldErrors.diskSizeGb}</p>
            ) : null}
          </div>

          <div className="field">
            <label htmlFor="rootVolumeType">rootVolumeType</label>
            <select
              id="rootVolumeType"
              value={form.rootVolumeType}
              onChange={(e) =>
                setForm({ ...form, rootVolumeType: e.target.value })
              }
            >
              {options.rootVolumeTypes.map((type) => (
                <option key={type} value={type}>
                  {type}
                </option>
              ))}
            </select>
          </div>

          <div className="field">
            <label htmlFor="dataDiskSizeGb">dataDiskSizeGb</label>
            <input
              id="dataDiskSizeGb"
              type="number"
              min={0}
              max={10000}
              value={form.dataDiskSizeGb}
              onChange={(e) =>
                setForm({ ...form, dataDiskSizeGb: Number(e.target.value) })
              }
            />
            {fieldErrors.dataDiskSizeGb ? (
              <p className="error">{fieldErrors.dataDiskSizeGb}</p>
            ) : null}
          </div>
        </div>

        <div className="checks">
          <label>
            <input
              type="checkbox"
              checked={form.staticIp}
              onChange={(e) => setForm({ ...form, staticIp: e.target.checked })}
            />
            staticIp
          </label>
        </div>

        {error ? <p className="error">{error}</p> : null}
        {toast ? <p className="toast">{toast}</p> : null}

        <div className="actions">
          <button className="btn btn-primary" type="submit" disabled={busy}>
            {busy ? "Envoi…" : "Créer la demande"}
          </button>
          <Link className="btn btn-ghost" to="/" style={{ display: "inline-flex" }}>
            Retour catalog
          </Link>
        </div>
      </form>
    </>
  );
}
