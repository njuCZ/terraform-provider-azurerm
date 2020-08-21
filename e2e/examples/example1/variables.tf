variable "prefix" {
  description = "The prefix that will be attached to all resources deployed"
  type        = string
  default     = "terratest-example"
}

variable "location" {
  description = "The location that all resources are deployed"
  type        = string
  default     = "southeastasia"
}