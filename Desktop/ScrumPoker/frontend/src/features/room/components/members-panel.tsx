import { Crown } from "lucide-react";

import { Card } from "../../../components/ui/card";
import { cn } from "../../../lib/utils";
import type { MemberView } from "../types";

type Props = {
  members: MemberView[];
  meId: number;
};

export function MembersPanel({ members, meId }: Props) {
  return (
    <Card className="space-y-3 p-4">
      <h3 className="text-base font-bold">اعضا</h3>

      <div className="space-y-2">
        {members.map((member) => (
          <div
            key={member.id}
            className={cn(
              "flex items-center justify-between rounded-xl border border-white/10 bg-surface-2/70 px-3 py-2",
              member.id === meId && "border-brand-400/50"
            )}
          >
            <div>
              <div className="flex items-center gap-2">
                <span className="font-semibold">{member.displayName}</span>
                {member.id === meId ? (
                  <span className="rounded-full bg-white/10 px-2 py-0.5 text-xs text-app-muted">شما</span>
                ) : null}
              </div>
              <p className="text-xs text-app-muted">{member.isActive ? "فعال" : "خارج"}</p>
            </div>

            {member.isHost ? (
              <span className="inline-flex items-center gap-1 rounded-full bg-brand-500/15 px-2 py-1 text-xs font-semibold text-brand-200">
                <Crown size={12} />
                میزبان
              </span>
            ) : null}
          </div>
        ))}
      </div>
    </Card>
  );
}
