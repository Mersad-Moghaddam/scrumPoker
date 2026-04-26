import { Card } from "../../../components/ui/card";
import { formatDateTime } from "../../../lib/utils";
import type { HistoryTask } from "../types";

type Props = {
  items: HistoryTask[];
};

export function HistoryPanel({ items }: Props) {
  return (
    <Card className="space-y-3 p-4">
      <h3 className="text-base font-bold">تاریخچه</h3>

      {items.length === 0 ? (
        <div className="rounded-xl border border-dashed border-white/10 bg-surface-2/70 p-3 text-sm text-app-muted">خالی</div>
      ) : (
        <div className="space-y-3">
          {items.map((item) => (
            <div key={item.id} className="rounded-xl border border-white/10 bg-surface-2/70 p-3">
              <div className="mb-2 flex items-start justify-between gap-3">
                <div>
                  <h4 className="font-bold">{item.title}</h4>
                  <p className="text-xs text-app-muted">{formatDateTime(item.completedAt ?? item.createdAt)}</p>
                </div>
                <span className="rounded-full bg-brand-500/15 px-2 py-1 text-xs font-semibold text-brand-200">
                  میانگین {item.averageLabel ?? "-"}
                </span>
              </div>
              <div className="grid gap-1.5">
                {item.votes.map((vote) => (
                  <div
                    key={`${item.id}-${vote.memberId}`}
                    className="flex items-center justify-between rounded-lg bg-white/5 px-2 py-1.5 text-sm"
                  >
                    <span>{vote.displayName}</span>
                    <span className="font-mono text-brand-100">{vote.value}</span>
                  </div>
                ))}
              </div>
            </div>
          ))}
        </div>
      )}
    </Card>
  );
}
