export function formatDate(value: string, withTime = false) {
  if (!value) return 'Não informado'
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return new Intl.DateTimeFormat('pt-BR', {
    dateStyle: 'medium',
    ...(withTime ? { timeStyle: 'short' as const } : {}),
  }).format(date)
}

export function toInputDate(value: string) {
  if (!value) return ''
  const date = new Date(value)
  const offset = date.getTimezoneOffset() * 60000
  return new Date(date.getTime() - offset).toISOString().slice(0, 16)
}

export function fromInputDate(value: string) { return value ? new Date(value).toISOString() : '' }
export function errorMessage(error: unknown, fallback: string) { return error instanceof Error ? error.message : fallback }
