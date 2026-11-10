import type { RequestStatus } from "../lib/provisioning";

const LABELS: Record<RequestStatus, string> = {
  submitted: "soumis",
  provisioning: "provisioning",
  ready: "ready",
  failed: "échec",
};

export function StatusBadge({ status }: { status: RequestStatus }) {
  return <span className={`badge ${status}`}>{LABELS[status]}</span>;
}
