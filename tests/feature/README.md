# ✔ feature/

* Controller → UseCase → Repository まで通す “結合テスト（Integration）”
* DB を実際に叩く or テスト用のスタブ DB を動かす
* ルーター設定やミドルウェアも含む

## 例：
* 「POST /users」を叩いて DB にレコードが作られること
* 「GET /users/{id}」で正しい JSON が返ること
* locale middleware が動いて翻訳されたレスポンスが返ること

＝ **API のふるまい全体** を保証するテスト

### 補足
* E2E ではないのでブラウザは不要
* Gin + httptest + MySQL（Docker）で十分
* CI中で実行する際はSQLiteとか用意しないとダメかも(着手時に要検討)
