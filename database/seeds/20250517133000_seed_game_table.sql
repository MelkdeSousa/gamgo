-- +goose Up
-- +goose StatementBegin
INSERT INTO games (
        id,
        title,
        description,
        platforms,
        releaseDate,
        rating,
        coverImage
    )
VALUES (
        '7038ad2e-2aef-47b6-9643-84589180b47e',
        'Halo 5',
        'Peace is interrupted when the colonies are unexpectedly attacked. When the greatest hero in the galaxy disappears, Spartan Locke is asked to locate Master Chief and solve the mystery that threatens the entire galaxy.',
        ARRAY ['PC'],
        '2021-05-27',
        0,
        'https://media.rawg.io/media/screenshots/1df/1df08ff4e8e40b15e0c7b09d4e569e35.jpg'
    ),
    (
        '523ccb10-be47-4349-abd1-738d48d703a2',
        'Halo 2',
        'The Covenant alien race threatens to destroy all humankind, and the only thing standing in its way is Master Chief, a genetically enhanced supersoldier.Master Chief returns in Halo 2, which features new vehicles, weapons, environments, and more.This time, you can interact with your environment, wield two weapons at the same time, board opponents vehicles, and even switch sides to play the role of a Covenant Elite. Halo 2 also supports broadband multiplayer action via Xbox Live.',
        ARRAY ['PC', 'Xbox'],
        '2004-11-09',
        438,
        'https://media.rawg.io/media/games/3bf/3bfc3bd9fda76bf83f6cf1d788e1c7c7.jpg'
    ),
    (
        'a1643215-4319-49bb-9dbf-883403f2205a',
        'Grand Theft Auto V',
        'Rockstar Games went bigger, since their previous installment of the series. You get the complicated and realistic world-building from Liberty City of GTA4 in the setting of lively and diverse Los Santos, from an old fan favorite GTA San Andreas.561 different vehicles (including every transport you can operate) and the amount is rising with every update. \nSimultaneous storytelling from three unique perspectives: \nFollow Michael, ex-criminal living his life of leisure away from the past, Franklin, a kid that seeks the better future, and Trevor, the exact past Michael is trying to run away from. \nGTA Online will provide a lot of additional challenge even for the experienced players, coming fresh from the story mode.Now you will have other players around that can help you just as likely as ruin your mission. Every GTA mechanic up to date can be experienced by players through the unique customizable character, and community content paired with the leveling system tends to keep everyone busy and engaged.',
        ARRAY ['PC', 'Playstation', 'Xbox'],
        '2013-09-17',
        447,
        'https://media.rawg.io/media/games/20a/20aa03a10cda45239fe22d035c0ebe64.jpg'
    );
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DELETE FROM games
WHERE id IN (
        '7038ad2e-2aef-47b6-9643-84589180b47e',
        '523ccb10-be47-4349-abd1-738d48d703a2',
        'a1643215-4319-49bb-9dbf-883403f2205a'
    );
-- +goose StatementEnd