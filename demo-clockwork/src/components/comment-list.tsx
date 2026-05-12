// comment-list.tsx
// Sigil custom component — minimal comment thread placeholder.
// Renderer emits flat props (datasource, taskField); both ignored in this
// demo-grade implementation. Real impl can replace the body when mock data
// wiring exists.

export interface CommentListProps {
  datasource?: string;
  taskField?: string;
  comments?: Array<Record<string, unknown>>;
  showAuthor?: boolean;
}

export function CommentList({ comments = [] }: CommentListProps) {
  if (comments.length === 0) {
    return <p className="text-sm text-muted-foreground">No comments yet.</p>;
  }
  return (
    <ul className="space-y-2 text-sm">
      {comments.map((c, i) => (
        <li key={i} className="rounded-md border p-2">
          <p className="text-xs text-muted-foreground">
            {String(c.author ?? "")} · {String(c.created_at ?? "")}
          </p>
          <p className="mt-1">{String(c.body ?? c.text ?? "")}</p>
        </li>
      ))}
    </ul>
  );
}

export default CommentList;
