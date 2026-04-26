import type { HistoryTask } from "../types";
import { formatDateTime } from "../../../lib/utils";
import { Card } from "../../../components/ui/card";

type Props = {
  items: HistoryTask[];
};

export function HistoryPanel({ items }: Props) {
  return (
    <Card className="space-y-4">
      <div>
        <h3 className="text-lg font-bold">تاریخچه جلسه</h3>
        <p className="text-sm text-app-muted">برآوردهای تکمیل‌شده با رأی هر عضو در اینجا می‌مانند.</p>
      </div>

      {items.length === 0 ? (
        <div className="rounded-3xl border border-dashed border-white/10 bg-surface-2/70 p-5 text-sm text-app-muted">
          هنوز هیچ تسکی کامل نشده است.
        </div>
      ) : (
        <div className="space-y-4">
          {items.map((item) => (
            <div key={item.id} className="rounded-3xl border border-white/10 bg-surface-2/70 p-4">
              <div className="mb-3 flex items-start justify-between gap-3">
                <div>
                  <h4 className="font-bold">{item.title}</h4>
                  <p className="text-xs text-app-muted">
                    تکمیل: {formatDateTime(item.completedAt ?? item.createdAt)}
                  </p>
                </div>
                <span className="rounded-full bg-brand-500/15 px-3 py-1 text-sm font-semibold text-brand-200">
                  میانگین: {item.averageLabel ?? "بدون برآورد عددی"}
                </span>
              </div>
              <div className="grid gap-2">
                {item.votes.map((vote) => (
                  <div
                    key={`${item.id}-${vote.memberId}`}
                    className="flex items-center justify-between rounded-2xl bg-white/5 px-3 py-2 text-sm"
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
