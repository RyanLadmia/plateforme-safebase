// Types pour les composants UI

export type ButtonVariant = 'primary' | 'secondary' | 'danger' | 'success' | 'warning'
export type ButtonSize = 'sm' | 'md' | 'lg'
export type ButtonType = 'button' | 'submit' | 'reset'

export interface ButtonProps {
  variant?: ButtonVariant
  size?: ButtonSize
  type?: ButtonType
  disabled?: boolean
  fullWidth?: boolean
  loading?: boolean
}

export type AlertType = 'success' | 'error' | 'warning' | 'info'

export interface AlertProps {
  type: AlertType
  title?: string
  message: string
  dismissible?: boolean
}

export interface ModalProps {
  isOpen: boolean
  title?: string
  size?: 'sm' | 'md' | 'lg' | 'xl'
  closable?: boolean
}

export interface FormFieldProps {
  label: string
  required?: boolean
  error?: string
  hint?: string
}
