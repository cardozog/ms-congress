type ConfirmModalProps = {
  title: string
  message: string
  confirmLabel: string
  tone?: 'publish' | 'close'
  loading?: boolean
  onConfirm: () => void
  onCancel: () => void
}

export function ConfirmModal({ title, message, confirmLabel, tone = 'publish', loading = false, onConfirm, onCancel }: ConfirmModalProps) {
  return <div className="modal-backdrop" role="presentation" onMouseDown={onCancel}><section className="confirm-modal" role="dialog" aria-modal="true" aria-labelledby="confirm-modal-title" onMouseDown={(event) => event.stopPropagation()}><span className={`modal-mark is-${tone}`}>{tone === 'publish' ? '↑' : '×'}</span><h2 id="confirm-modal-title">{title}</h2><p>{message}</p><div className="modal-actions"><button className="secondary-button" type="button" onClick={onCancel} disabled={loading}>Cancelar</button><button className={`modal-confirm-button is-${tone}`} type="button" onClick={onConfirm} disabled={loading}>{loading ? 'Aguarde...' : confirmLabel}</button></div></section></div>
}
