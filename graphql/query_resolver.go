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

	if pagination != nil {
		res, err := r.server.accountClient.GetAccounts(ctx, uint64(*pagination.Skip), uint64(*pagination.Take))
		if err != nil {
			return nil, err
		}
		accounts := make([]*Account, len(res))
		for i, account := range res {
			accounts[i] = &Account{
				ID: account.ID,
				Name: account.Name,
			}
		}
		return accounts, nil
	}
	return nil, nil
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

	if pagination != nil {
		res, err := r.server.catalogClient.GetProducts(ctx, uint64(*pagination.Skip), uint64(*pagination.Take))
		if err != nil {
			return nil, err
		}
		products := make([]*Product, len(res))
		for i, p := range res {
			products[i] = &Product{
				ID: p.ID,
				Name: p.Name,
				Price: p.Price,
				Description: p.Description,
			}
		}
		return products, nil
	}

	if query != nil {

		res, err := r.server.catalogClient.SearchProducts(ctx, *query, uint64(*pagination.Skip), uint64(*pagination.Take))
		if err != nil {
			return nil, err
		}
		products := make([]*Product, len(res))
		for i, p := range res {
			products[i] = &Product{
				ID: p.ID,
				Name: p.Name,
				Price: p.Price,
				Description: p.Description,
			}
		}
		return products, nil
	}
	return nil, nil
}


