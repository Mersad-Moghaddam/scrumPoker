import { FormEvent, useState } from "react";

import { Button } from "../../../components/ui/button";
import { Card } from "../../../components/ui/card";

type Props = {
  onSubmit: (title: string) => void;
  pending?: boolean;
  title?: string;
};

export function CreateTaskCard({ onSubmit, pending, title = "تسک" }: Props) {
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
    <Card className="space-y-3 p-4">
      <h3 className="text-base font-bold">{title}</h3>
      <form onSubmit={handleSubmit} className="space-y-3">
        <input
          value={taskTitle}
          onChange={(event) => setTaskTitle(event.target.value)}
          className="w-full rounded-xl border border-white/10 bg-surface-2 px-3 py-2 outline-none transition focus:border-brand-400"
          placeholder="عنوان"
        />
        <Button type="submit" className="w-full" disabled={!taskTitle.trim() || pending}>
          {pending ? "..." : "شروع"}
        </Button>
      </form>
    </Card>
  );
}
