-- 以下のような初期データがある場合は記述する
INSERT IGNORE INTO notification_types
    (`id`, `title`)
VALUES
    (1, "ポイント送付のお知らせ"),
    (2, "その他"),
    (3, "キャンペーン");
