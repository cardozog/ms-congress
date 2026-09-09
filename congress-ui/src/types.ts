export interface Session {
  token: string
  tokenType?: string
  email: string
  id: number
  tipoPessoaID: number
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

export interface EventDetails extends EventSummary {
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
  tipoIngresso: { id: number; nome: string }
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
