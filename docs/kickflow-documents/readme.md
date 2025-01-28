# Kickflow Terraform Provider の設計

## リソース (Resource)

### 1. ユーザー関連

#### `kickflow_user`
- ユーザーを管理するリソース
- 属性:
  - email (必須): メールアドレス
  - code: ユーザーコード（未指定の場合自動生成）
  - first_name (必須): 名
  - last_name (必須): 姓
  - employee_id: 社員番号
  - send_email: 招待メールを送信するかどうか（デフォルト: true）

#### `kickflow_grade`
- 役職を管理するリソース
- 属性:
  - name (必須): 役職名
  - level (必須): レベル（0-255）
  - code: 役職コード（未指定の場合自動生成）
  - is_default: デフォルトの役職かどうか

#### `kickflow_proxy_applicant`
- 代理申請の設定を管理するリソース
- 属性:
  - user_id (必須): 代理される側のユーザーID
  - proxy_user_id (必須): 代理する側のユーザーID
  - starts_on: 開始日（未指定の場合即時開始）
  - ends_on: 終了日（未指定の場合無期限）
  - workflow_ids: 対象ワークフローのID配列

#### `kickflow_proxy_approver`
- 代理承認の設定を管理するリソース
- 属性:
  - user_id (必須): 代理される側のユーザーID
  - proxy_user_id (必須): 代理する側のユーザーID
  - starts_on: 開始日（未指定の場合即時開始）
  - ends_on: 終了日（未指定の場合無期限）
  - workflow_ids: 対象ワークフローのID配列

### 2. 組織関連

#### `kickflow_organization_chart`
- 組織図を管理するリソース
- 属性:
  - name (必須): 組織図名
  - current: 現在有効な組織図かどうか
  - due_on: 有効化予定日

#### `kickflow_team`
- チームを管理するリソース
- 属性:
  - name (必須): チーム名
  - organization_chart_id (必須): 所属する組織図のID
  - code: チームコード（未指定の場合自動生成）
  - parent_id: 親チームのID
  - approve_only: 承認専用チームかどうか

#### `kickflow_team_membership`
- チームメンバーシップを管理するリソース
- 属性:
  - team_id (必須): チームID
  - user_id (必須): ユーザーID
  - leader: 上長かどうか
  - grade_ids: 役職ID配列

### 3. ワークフロー関連

#### `kickflow_category`
- カテゴリを管理するリソース
- 属性:
  - name (必須): カテゴリ名

#### `kickflow_folder`
- フォルダを管理するリソース
- 属性:
  - name (必須): フォルダ名
  - code: フォルダコード（未指定の場合自動生成）
  - description: 説明
  - parent_id: 親フォルダのID

### 4. その他

#### `kickflow_general_master`
- 汎用マスタを管理するリソース
- 属性:
  - name (必須): 汎用マスタ名
  - code: コード（未指定の場合自動生成）
  - description: 説明
  - fields: カスタムフィールド定義の配列

#### `kickflow_general_master_item`
- 汎用マスタアイテムを管理するリソース
- 属性:
  - general_master_id (必須): 汎用マスタID
  - name (必須): アイテム名
  - code: コード（未指定の場合自動生成）
  - description: 説明
  - starts_on: 有効期限開始日
  - ends_on: 有効期限終了日
  - inputs: カスタムフィールドの入力値

#### `kickflow_role`
- 管理者ロールを管理するリソース
- 属性:
  - name (必須): ロール名
  - permission_list: 権限リスト

## データソース (Data Source)

### 1. ユーザー関連

#### `kickflow_user`
- ユーザー情報を取得
- 検索条件:
  - id: ユーザーID
  - email: メールアドレス
  - code: ユーザーコード
- 取得データ:
  - id: UUID
  - email: メールアドレス
  - code: ユーザーコード
  - first_name: 名
  - last_name: 姓
  - full_name: フルネーム
  - employee_id: 社員番号
  - image: ユーザー画像（複数サイズ）
  - status: ステータス（invited/activated/suspended/deactivated）
  - locale: 言語設定
  - created_at: 作成日時
  - updated_at: 更新日時
  - deactivated_at: 削除日時

#### `kickflow_users`
- ユーザー一覧を取得
- フィルタ条件:
  - status: ステータス（invited/activated/suspended/deactivated）
  - page: ページ番号（1から開始）
  - per_page: 1ページあたりの件数（最大100）
- 取得データ:
  - users: ユーザー情報の配列（userと同じ構造）

#### `kickflow_grade`
- 役職情報を取得
- 検索条件:
  - id: 役職ID
  - code: 役職コード
- 取得データ:
  - id: UUID
  - name: 役職名
  - level: レベル（0-255）
  - code: 役職コード
  - is_default: デフォルトの役職かどうか
  - created_at: 作成日時
  - updated_at: 更新日時

#### `kickflow_grades`
- 役職一覧を取得
- フィルタ条件:
  - page: ページ番号（1から開始）
  - per_page: 1ページあたりの件数（最大100）
  - sort_by: ソート項目（level/code）
- 取得データ:
  - grades: 役職情報の配列（gradeと同じ構造）

### 2. 組織関連

#### `kickflow_organization_chart`
- 組織図情報を取得
- 検索条件:
  - id: 組織図ID
  - current: 現在有効な組織図のみ取得
- 取得データ:
  - id: UUID
  - name: 組織図名
  - current: 現在有効な組織図かどうか
  - teams_count: チーム数
  - memberships_count: 所属数
  - created_at: 作成日時
  - updated_at: 更新日時
  - activation_plan: 有効化予定情報

#### `kickflow_organization_charts`
- 組織図一覧を取得
- フィルタ条件:
  - page: ページ番号（1から開始）
  - per_page: 1ページあたりの件数（最大100）
  - sort_by: ソート項目（created_at/name）
- 取得データ:
  - organization_charts: 組織図情報の配列（organization_chartと同じ構造）

#### `kickflow_team`
- チーム情報を取得
- 検索条件:
  - id: チームID
  - code: チームコード
- 取得データ:
  - id: UUID
  - name: チーム名
  - full_name: 上位組織を含む名前
  - code: チームコード
  - approve_only: 承認専用チームかどうか
  - users_count: ユーザー数
  - created_at: 作成日時
  - updated_at: 更新日時
  - parent: 親チーム情報
  - children: 子チーム情報の配列
  - users: メンバー情報の配列

#### `kickflow_teams`
- チーム一覧を取得
- フィルタ条件:
  - organization_chart_id: 組織図ID
  - parent_id: 親チームID
  - page: ページ番号（1から開始）
  - per_page: 1ページあたりの件数（最大100）
- 取得データ:
  - teams: チーム情報の配列（teamと同じ構造）

### 3. ワークフロー関連

#### `kickflow_category`
- カテゴリ情報を取得
- 検索条件:
  - id: カテゴリID
- 取得データ:
  - id: UUID
  - name: カテゴリ名
  - created_at: 作成日時
  - updated_at: 更新日時

#### `kickflow_categories`
- カテゴリ一覧を取得
- フィルタ条件:
  - page: ページ番号（1から開始）
  - per_page: 1ページあたりの件数（最大100）
  - sort_by: ソート項目（name/created_at/updated_at）
- 取得データ:
  - categories: カテゴリ情報の配列（categoryと同じ構造）

#### `kickflow_folder`
- フォルダ情報を取得
- 検索条件:
  - id: フォルダID
  - code: フォルダコード
- 取得データ:
  - id: UUID
  - name: フォルダ名
  - code: フォルダコード
  - description: 説明
  - workflows_count: ワークフロー数
  - routes_count: 経路数
  - pipelines_count: パイプライン数
  - created_at: 作成日時
  - updated_at: 更新日時
  - ancestors: 親フォルダ情報の配列
  - children: 子フォルダ情報の配列

#### `kickflow_folders`
- フォルダ一覧を取得
- フィルタ条件:
  - page: ページ番号（1から開始）
  - per_page: 1ページあたりの件数（最大100）
  - sort_by: ソート項目（created_at/name）
- 取得データ:
  - folders: フォルダ情報の配列（folderと同じ構造）

#### `kickflow_workflow`
- ワークフロー情報を取得
- 検索条件:
  - id: ワークフローID
  - code: ワークフローコード
- 取得データ:
  - id: UUID
  - code: ワークフローコード
  - version_id: バージョンID
  - version_number: バージョン番号
  - name: ワークフロー名
  - description: 説明
  - status: ステータス（visible/invisible/deleted）
  - public_ticket: チケットが全体公開されるかどうか
  - visible_to_manager: 申請者の上長への公開設定
  - visible_to_team_members: 申請チームのメンバーへの公開設定
  - title_description: タイトルの説明
  - ticket_number_format: チケット番号フォーマット
  - overwritable: 承認者による上書き可否
  - created_at: 作成日時
  - updated_at: 更新日時
  - title_input_mode: タイトル入力モード
  - title_formula: タイトルの計算式
  - allow_editing_of_viewers: 閲覧者編集可否
  - author: 作成者情報
  - version_author: バージョン作成者情報
  - folder: フォルダ情報
  - categories: カテゴリ情報の配列

#### `kickflow_workflows`
- ワークフロー一覧を取得
- フィルタ条件:
  - status: ステータス（visible/invisible）
  - page: ページ番号（1から開始）
  - per_page: 1ページあたりの件数（最大100）
  - sort_by: ソート項目（created_at/updated_at/name/status）
- 取得データ:
  - workflows: ワークフロー情報の配列（workflowと同じ構造）

### 4. その他

#### `kickflow_general_master`
- 汎用マスタ情報を取得
- 検索条件:
  - id: 汎用マスタID
  - code: 汎用マスタコード
- 取得データ:
  - id: UUID
  - code: コード
  - name: 汎用マスタ名
  - description: 説明
  - created_at: 作成日時
  - updated_at: 更新日時
  - fields: カスタムフィールド定義の配列

#### `kickflow_general_masters`
- 汎用マスタ一覧を取得
- フィルタ条件:
  - page: ページ番号（1から開始）
  - per_page: 1ページあたりの件数（最大100）
  - sort_by: ソート項目（created_at/code/name）
- 取得データ:
  - general_masters: 汎用マスタ情報の配列（general_masterと同じ構造）

#### `kickflow_general_master_item`
- 汎用マスタアイテム情報を取得
- 検索条件:
  - id: アイテムID
  - code: アイテムコード
  - general_master_id: 汎用マスタID
- 取得データ:
  - id: UUID
  - code: コード
  - name: アイテム名
  - description: 説明
  - created_at: 作成日時
  - updated_at: 更新日時
  - starts_on: 有効期限開始日
  - ends_on: 有効期限終了日
  - inputs: カスタムフィールドの入力値配列

#### `kickflow_general_master_items`
- 汎用マスタアイテム一覧を取得
- フィルタ条件:
  - general_master_id: 汎用マスタID
  - page: ページ番号（1から開始）
  - per_page: 1ページあたりの件数（最大100）
  - sort_by: ソート項目（created_at/code/name）
- 取得データ:
  - general_master_items: 汎用マスタアイテム情報の配列（general_master_itemと同じ構造）

#### `kickflow_role`
- 管理者ロール情報を取得
- 検索条件:
  - id: ロールID
- 取得データ:
  - id: UUID
  - name: ロール名
  - editable: 編集可能かどうか
  - users_count: 所属ユーザー数
  - created_at: 作成日時
  - updated_at: 更新日時
  - permission_list: 権限リスト（詳細な権限設定情報）

#### `kickflow_roles`
- 管理者ロール一覧を取得
- フィルタ条件:
  - page: ページ番号（1から開始）
  - per_page: 1ページあたりの件数（最大100）
  - sort_by: ソート項目（created_at/name）
- 取得データ:
  - roles: 管理者ロール情報の配列（roleと同じ構造）
