package post

import "database/sql"

type Post struct {
	ID      int
	Content string
	Author  string
	DB      *sql.DB
}

func (p *Post) Create() error {
	_, err := p.DB.Exec("insert into posts (id, content, author) values (?, ?, ?)", p.ID, p.Content, p.Author)
	return err
}

func (p *Post) GetPost(id int) (Post, error) {
	row := p.DB.QueryRow("select id, content, author from posts where id = ?", id)
	err := row.Scan(&p.ID, &p.Content, &p.Author)
	return *p, err
}
