// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mq_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/mq"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccMQBrokerEngineTypesDataSource_basic(t *testing.T) {
	ctx := acctest.Context(t)
	dataSourceName := "data.aws_mq_broker_engine_types.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, mq.EndpointsID)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, mq.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccBrokerEngineTypesDataSourceConfig_basic("ACTIVEMQ"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(dataSourceName, "broker_engine_types.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "broker_engine_types.0.engine_type", "ACTIVEMQ"),
				),
			},
		},
	})
}

func testAccBrokerEngineTypesDataSourceConfig_basic(engineType string) string {
	return fmt.Sprintf(`
data "aws_mq_broker_engine_types" "test" {
  engine_type = %[1]q
}
`, engineType)
}
