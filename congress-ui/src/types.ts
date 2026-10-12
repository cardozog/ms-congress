export interface Session {
  token: string
  tokenType?: string
  email: string
  id: number
  tipoPessoaID: number
  papeis: string[]
}

export interface CompanyProfile {
  id?: number
  usuarioID?: number
  razaoSocial?: string
  nomeFantasia?: string
}

export interface Profile {
  id: number
  email: string
  pessoaJuridica?: CompanyProfile | null
}

export interface EventSummary {
  id: number
  nome: string
  descricao: string
  dataInicio: string
  dataFim: string
  publicado: boolean
  encerrado: boolean
  organizadorID: number
  eventoLogo?: string
}

export interface Organizer {
  id: number
  usuarioID: number
  razaoSocial: string
  nomeFantasia: string
  cnpj: string
}

export interface AvailableEvent extends EventSummary {
  organizador: Organizer
}

export interface AvailableEventGroup {
  organizador: Organizer
  eventos: AvailableEvent[]
}

export interface EventDetails extends EventSummary {
  organizador: Organizer
  logradouro: string
  cidade: string
  estado: string
  cep: string
  ingressos: TicketDetails[]
}

export interface TicketDetails {
  id: number
  eventoID: number
  tipoIngressoID: number
  preco: number
  quantidade: number
  disponivel: number
  dataInicioVenda: string
  dataFimVenda: string
  tipoIngresso: { id: number; nome: string }
}

export interface CartItem {
  id: number
  evento: { id: number; nome: string; eventoLogo?: string }
  ingresso: { id: number; nome: string }
  quantidade: number
  precoUnitario: number
  subtotal: number
}

export interface Cart {
  itens: CartItem[]
  totalItens: number
  valorTotal: number
}

export interface TicketPayload {
  tipo_ingresso_id: number
  preco: number
  quantidade: number
}

export interface TicketType {
  id: number
  nome: string
}

export interface EventPayload {
  nome: string
  descricao: string
  data_inicio: string
  data_fim: string
  organizador_id: number
  logradouro: string
  cidade: string
  estado: string
  cep: string
  ingressos: TicketPayload[]
}
