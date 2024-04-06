package qrcierrelotesdtos

import "github.com/Corrientes-Telecomunicaciones/api_go_pasarela_checkout/pkg/entities"

type RequestCreateQrCierrelotes struct {
	UpdatePagointento bool                  `json:"update_pagointento"`
	UpdatePago        bool                  `json:"update_pago"`
	QrCierrelote      entities.Qrcierrelote `json:"qrcierrelote"`
	Pagointento       entities.Pagointento  `json:"pagointento"`
	Pago              entities.Pago         `json:"pago"`
}
