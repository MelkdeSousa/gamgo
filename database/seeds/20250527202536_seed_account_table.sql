-- +goose Up
-- +goose StatementBegin
-- valid account
INSERT INTO accounts (
        id,
        name,
        email,
        passwordHash,
        isActive,
        createdAt
    )
VALUES (
        '36582056-2e6f-4b98-b267-5ed080d9e30b',
        'John Doe',
        'john.doe@example.com',
        '$2a$14$sRnkgTyK2Fhc9SXgiS8cpODV5q5HZrz38xyHMPARroYEefWES4RVG', -- secret
        true,
        NOW()
    ),
    (
        '97ccaf95-00d7-4493-9003-d6d0bfd8a81b',
        'Alice Johnson',
        'alice.johnson@example.com',
        '$2a$14$ew038CG9rzOdqpiLMX63X.MpmBfrNjRE0umfH36RdDtKYdEmbIvY6', -- turtle
        false,
        NOW()
    );
-- deleted account
INSERT INTO accounts (
        id,
        name,
        email,
        passwordHash,
        isActive,
        createdAt,
        deletedAt
    )
VALUES (
        '36725e1f-2ff7-4f53-a86c-553f1ea30eb5',
        'Jane Smith',
        'jane.smith@example.com',
        '$2a$14$zqeTEaYTMyroE.xzZIk8q.yHFzcG8jVGt2df8zZeWowWCYpLd5xhK', -- burble
        false,
        NOW(),
        NOW()
    );
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DELETE FROM accounts
WHERE id IN (
        '36582056-2e6f-4b98-b267-5ed080d9e30b',
        '36725e1f-2ff7-4f53-a86c-553f1ea30eb5',
        '97ccaf95-00d7-4493-9003-d6d0bfd8a81b'
    );
-- +goose StatementEnd