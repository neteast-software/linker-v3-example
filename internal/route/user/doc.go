// Package user declares user HTTP routes.
//
// Route is the HTTP entry boundary, not a thin forwarding layer. A route file
// should keep the route declaration, resource scope, request payload, header
// parsing, response shape and short API-only flow close together. Reusable
// business capability belongs in service.
package user
