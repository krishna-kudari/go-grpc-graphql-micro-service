package main

import (
	"context"
	"log"
	"time"
)

type queryResolver struct {
	server *Server
}

func (r *queryResolver) Accounts(ctx context.Context, pagination *PaginationInput, id *string)([]*Account, error)  {
	context, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if id != nil{
		r, err := r.server.accountClient.GetAccountInfo(context, *id)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return []*Account{{
			ID: r.ID,
			Name: r.Name,
		}}, nil
	}

	// Default pagination when not provided
	skip, take := uint64(0), uint64(100)
	if pagination != nil {
		if pagination.Skip != nil {
			skip = uint64(*pagination.Skip)
		}
		if pagination.Take != nil && *pagination.Take > 0 {
			take = uint64(*pagination.Take)
		}
	}
	res, err := r.server.accountClient.GetAccounts(ctx, skip, take)
	if err != nil {
		return nil, err
	}
	accounts := make([]*Account, len(res))
	for i, account := range res {
		accounts[i] = &Account{
			ID:   account.ID,
			Name: account.Name,
		}
	}
	return accounts, nil
}

func (r *queryResolver) Products(ctx context.Context, pagination *PaginationInput, query *string, id *string)([]*Product, error)  {
	context, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if id != nil {
		r, err := r.server.catalogClient.GetProduct(context, *id)
		if err != nil {
			return nil, err
		}
		return []*Product{{
			ID: r.ID,
			Name: r.Name,
			Price: r.Price,
			Description: r.Description,
		}},nil
	}

	if query != nil {
		skip, take := uint64(0), uint64(100)
		if pagination != nil {
			if pagination.Skip != nil {
				skip = uint64(*pagination.Skip)
			}
			if pagination.Take != nil && *pagination.Take > 0 {
				take = uint64(*pagination.Take)
			}
		}
		res, err := r.server.catalogClient.SearchProducts(ctx, *query, skip, take)
		if err != nil {
			return nil, err
		}
		products := make([]*Product, len(res))
		for i, p := range res {
			products[i] = &Product{
				ID:          p.ID,
				Name:        p.Name,
				Price:       p.Price,
				Description: p.Description,
			}
		}
		return products, nil
	}

	// Default pagination when not provided
	skip, take := uint64(0), uint64(100)
	if pagination != nil {
		if pagination.Skip != nil {
			skip = uint64(*pagination.Skip)
		}
		if pagination.Take != nil && *pagination.Take > 0 {
			take = uint64(*pagination.Take)
		}
	}
	res, err := r.server.catalogClient.GetProducts(ctx, skip, take)
	if err != nil {
		return nil, err
	}
	products := make([]*Product, len(res))
	for i, p := range res {
		products[i] = &Product{
			ID:          p.ID,
			Name:        p.Name,
			Price:       p.Price,
			Description: p.Description,
		}
	}
	return products, nil
}


