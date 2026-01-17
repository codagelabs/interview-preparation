# Parent-Child Archival Policy System - Problem Statement

## Executive Summary

The current Parent-Child Archival Policy system suffers from a critical database table explosion issue that creates significant operational overhead, performance bottlenecks, and unmanageable schema growth. This document outlines the problem, its impact, and the need for architectural improvements.

## System Overview

### Application Context
- **Purpose**: Archival application that manages data lifecycle based on configurable policies
- **Architecture**: Hierarchical policy system with parent-child relationships
- **Execution Model**: Cron-based scheduling with session-based processing

### Policy Configuration
Each archival policy contains:
- **Name**: Unique identifier for the policy
- **Cron Expression**: Defines when the archival job executes
- **Callback URL**: Endpoint invoked during cron execution
- **Archival Criteria**: Conditions determining which entities should be archived

### Session Model
- **Definition**: One session = one cron execution of a policy
- **Scope**: During a session, entities matching criteria are archived
- **Lifecycle**: Session-based processing with entity ID tracking

## Current Design

### Parent-Child Policy Relationship
Policies are organized in hierarchical relationships:
- **Parent Policy**: Archives primary entities (e.g., orders)
- **Child Policy**: Processes archived entity IDs from parent (e.g., deliveries linked to orders)
- **Cascading**: Child policies can forward IDs to their own children

### Data Flow Example
```
Parent Policy (Order Service)
    ↓ (archives orders, sends order IDs)
Child Policy (Delivery Management Service)
    ↓ (archives deliveries linked to those orders)
Grandchild Policy (if applicable)
```

### Current Implementation
For each parent session, when entity IDs are sent to child policies:

1. **Table Creation**: A new table is created per child policy per parent session
2. **Data Storage**: Archived entity IDs are stored in session-specific tables
3. **Child Processing**: Child policy cron reads from its respective table
4. **Cascading**: Child can forward IDs to its own child policies

## The Problem: Table Explosion

### Current Behavior Example
**Scenario**: 1 parent policy → 3 child policies (c1, c2, c3)
- **Sessions**: Parent runs 10 sessions
- **Table Creation**: Each session creates 3 child tables (c1, c2, c3)
- **Total Tables**: 10 × 3 = 30 tables

**Scaled Scenario**: 1 parent policy → 2 child policies
- **Sessions**: Parent runs 10 sessions  
- **Total Tables**: 10 × 2 = 20 tables

### Mathematical Impact
```
Total Tables = Number of Parent Sessions × Number of Child Policies
```

**Examples**:
- 5 parent policies × 3 children × 20 sessions = 300 tables
- 10 parent policies × 5 children × 50 sessions = 2,500 tables
- 20 parent policies × 10 children × 100 sessions = 20,000 tables

## Impact Analysis

### 1. Database Management Overhead
- **Schema Proliferation**: Exponential growth of table count
- **Maintenance Complexity**: Difficult to track and manage thousands of tables
- **Backup Challenges**: Increased backup time and storage requirements
- **Migration Complexity**: Schema changes become increasingly difficult

### 2. Performance Bottlenecks
- **Query Performance**: Database catalog becomes bloated
- **Connection Overhead**: Increased metadata operations
- **Index Management**: Proliferation of indexes across many tables
- **Memory Usage**: Database metadata consumes excessive memory

### 3. Operational Challenges
- **Monitoring Difficulty**: Hard to track table usage and performance
- **Cleanup Complexity**: Identifying and removing obsolete tables
- **Resource Allocation**: Unpredictable storage and compute requirements
- **Scaling Limitations**: System becomes harder to scale horizontally

### 4. Development and Maintenance Issues
- **Code Complexity**: Dynamic table creation logic becomes complex
- **Testing Challenges**: Difficult to test with varying table structures
- **Documentation**: Hard to maintain accurate schema documentation
- **Debugging**: Troubleshooting becomes increasingly difficult

## Root Cause Analysis

### Primary Causes
1. **Session-Based Table Creation**: New table per session per child policy
2. **Lack of Data Consolidation**: No mechanism to reuse or consolidate tables
3. **Missing Lifecycle Management**: No automated cleanup of obsolete tables
4. **Inflexible Schema Design**: Hard-coded approach to data storage

### Contributing Factors
- **Short-Term Thinking**: Solution designed for immediate needs without scalability consideration
- **Lack of Data Partitioning Strategy**: No consideration for data organization
- **Missing Archival Strategy**: No plan for managing historical session data
- **Insufficient Monitoring**: No visibility into table growth patterns

## Business Impact

### Immediate Risks
- **System Degradation**: Performance issues affecting user experience
- **Operational Costs**: Increased infrastructure and maintenance costs
- **Development Velocity**: Slower feature development due to complexity
- **Reliability Concerns**: Higher risk of system failures

### Long-Term Consequences
- **Scalability Limits**: System cannot grow beyond certain thresholds
- **Technical Debt**: Accumulating complexity that becomes harder to address
- **Competitive Disadvantage**: Slower response to business requirements
- **Resource Waste**: Inefficient use of database resources

## Success Criteria for Solution

### Functional Requirements
- **Data Integrity**: Maintain all current functionality
- **Performance**: Equal or better performance than current system
- **Scalability**: Support 10x current scale without degradation
- **Backward Compatibility**: Minimal disruption to existing operations

### Non-Functional Requirements
- **Maintainability**: Easier to manage and monitor
- **Extensibility**: Support for future policy types and relationships
- **Reliability**: Improved system stability and error handling
- **Cost Efficiency**: Reduced infrastructure and operational costs

## Recommended Next Steps

### Phase 1: Immediate Actions
1. **Assessment**: Complete analysis of current table usage patterns
2. **Monitoring**: Implement comprehensive table growth monitoring
3. **Documentation**: Create detailed current state documentation
4. **Stakeholder Alignment**: Ensure all teams understand the problem scope

### Phase 2: Solution Design
1. **Architecture Review**: Design new data storage approach
2. **Proof of Concept**: Validate proposed solution with limited scope
3. **Migration Strategy**: Plan transition from current to new system
4. **Risk Assessment**: Identify and mitigate implementation risks

### Phase 3: Implementation
1. **Pilot Implementation**: Deploy solution with subset of policies
2. **Performance Testing**: Validate performance under load
3. **Gradual Migration**: Move policies to new system incrementally
4. **Monitoring and Optimization**: Continuous improvement based on metrics

## Conclusion

The current Parent-Child Archival Policy system's table explosion problem represents a critical scalability and maintainability issue that requires immediate attention. The exponential growth of database tables creates significant operational overhead and performance bottlenecks that will only worsen over time.

A comprehensive solution is needed that addresses the root causes while maintaining system functionality and improving scalability. The recommended approach involves architectural redesign with proper data consolidation, lifecycle management, and monitoring capabilities.

**Priority**: High - This issue should be addressed before it significantly impacts system performance and operational efficiency.

---

*Document Version: 1.0*  
*Last Updated: [Current Date]*  
*Author: System Architecture Team*


