// Command award-prs awards the merged_prs badges from config/rules.yaml to
// the authors of the merged pull requests stored in the github_* tables.
// It is safe to re-run: a person never gets the same badge twice.
package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"tahrir-go/internal/db"
	"tahrir-go/internal/models"
	"tahrir-go/internal/rules"
	worker "tahrir-go/internal/rules/workers"
)

func main() {
	_ = godotenv.Load()
	conn, err := db.Connect(db.DSNFromEnv())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	rulesFile := os.Getenv("RULES_FILE")
	if rulesFile == "" {
		rulesFile = "config/rules.yaml"
	}
	configuredRules, err := rules.LoadRules(rulesFile, conn)
	if err != nil {
		log.Fatalf("failed to load rules: %v", err)
	}

	// Only the merged-PR rules apply here; the others react to other events.
	var prRules []rules.Rule
	for _, rule := range configuredRules {
		prRule, ok := rule.(*rules.MergedPRsRule)
		if !ok {
			continue
		}
		if err := db.EnsureBadge(conn, models.Badge{
			ID:          prRule.BadgeID,
			Name:        prRule.Name,
			Description: prRule.Description,
			Criteria:    prRule.Description,
			CreatedOn:   time.Now(),
			Tags:        "github,pull-request",
		}); err != nil {
			log.Fatalf("failed to create badge %q: %v", prRule.BadgeID, err)
		}
		prRules = append(prRules, prRule)
	}
	if len(prRules) == 0 {
		log.Fatalf("no enabled merged_prs rules in %s", rulesFile)
	}

	linked, err := db.LinkMergedPRAuthors(conn)
	if err != nil {
		log.Fatalf("failed to link PR authors to persons: %v", err)
	}

	personIDs, err := db.PersonIDsWithMergedPRs(conn)
	if err != nil {
		log.Fatalf("failed to list PR authors: %v", err)
	}

	w := worker.Worker{DB: conn, Rules: prRules}
	awarded := 0
	for _, personID := range personIDs {
		awarded += w.Handle(rules.Event{Type: rules.PullRequestMerged, PersonID: personID})
	}

	log.Printf("linked %d new PR authors, checked %d people, awarded %d badges", linked, len(personIDs), awarded)
}
