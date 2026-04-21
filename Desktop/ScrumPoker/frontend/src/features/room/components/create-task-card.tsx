import { FormEvent, useState } from "react";

import { Button } from "../../../components/ui/button";
import { Card } from "../../../components/ui/card";

type Props = {
  onSubmit: (title: string) => void;
  pending?: boolean;
  title?: string;
  description?: string;
};

export function CreateTaskCard({
  onSubmit,
  pending,
  title = "تسک بعدی را شروع کن",
  description = "برای MVP فقط عنوان تسک کافی است."
}: Props) {
  const [taskTitle, setTaskTitle] = useState("");

  const handleSubmit = (event: FormEvent) => {
    event.preventDefault();
    if (!taskTitle.trim()) {
      return;
    }
    onSubmit(taskTitle.trim());
    setTaskTitle("");
  };

  return (
    <Card className="space-y-4">
      <div>
        <h3 className="text-lg font-bold">{title}</h3>
        <p className="text-sm text-app-muted">{description}</p>
      </div>
      <form onSubmit={handleSubmit} className="space-y-4">
        <input
          value={taskTitle}
          onChange={(event) => setTaskTitle(event.target.value)}
          className="w-full rounded-2xl border border-white/10 bg-surface-2 px-4 py-3 outline-none transition focus:border-brand-400"
          placeholder="مثلا طراحی API برای گزارش‌گیری"
        />
        <Button type="submit" className="w-full" disabled={!taskTitle.trim() || pending}>
          {pending ? "در حال شروع..." : "شروع رأی‌گیری"}
        </Button>
      </form>
    </Card>
  );
}
